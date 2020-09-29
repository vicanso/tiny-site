// Copyright 2019 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"net/http"
	"time"

	"github.com/dustin/go-humanize"
	warner "github.com/vicanso/count-warner"
	"github.com/vicanso/elton"
	eltonMid "github.com/vicanso/elton/middleware"
	"github.com/vicanso/hes"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"go.uber.org/zap"

	"github.com/vicanso/tiny-site/config"
	_ "github.com/vicanso/tiny-site/controller"
	_ "github.com/vicanso/tiny-site/helper"
	"github.com/vicanso/tiny-site/log"
	"github.com/vicanso/tiny-site/middleware"
	"github.com/vicanso/tiny-site/router"
	_ "github.com/vicanso/tiny-site/schedule"
	"github.com/vicanso/tiny-site/service"
	"github.com/vicanso/tiny-site/util"
)

// 相关依赖服务的校验，主要是数据库等
func dependServiceCheck() (err error) {
	err = service.RedisPing()
	if err != nil {
		return
	}
	configSrv := new(service.ConfigurationSrv)
	err = configSrv.Refresh()
	if err != nil {
		return
	}
	return
}

func main() {
	logger := log.Default()
	e := elton.New()
	e.SignedKeys = service.GetSignedKeys()

	// 未处理的error才会触发
	// 如果1分钟出现超过5次未处理异常
	warnerException := warner.NewWarner(60*time.Second, 5)
	warnerException.ResetOnWarn = true
	warnerException.On(func(_ string, _ warner.Count) {
		service.AlarmError("too many uncaught exception")
	})
	e.OnError(func(c *elton.Context, err error) {
		if !util.IsProduction() {
			he, ok := err.(*hes.Error)
			if ok {
				if he.Extra == nil {
					he.Extra = make(map[string]interface{})
				}
				he.Extra["stack"] = util.GetStack(5)
			}
		}

		// 可以针对实际场景输出更多的日志信息
		logger.DPanic("exception",
			zap.String("ip", c.RealIP()),
			zap.String("uri", c.Request.RequestURI),
			zap.Error(err),
		)
		// TODO 邮件通知
		warnerException.Inc("exception", 1)
	})
	// 对于404的请求，不会执行中间件，一般都是因为攻击之类才会导致大量出现404，
	// 因此可在此处汇总出错IP，针对较频繁出错IP，增加告警信息
	// 如果1分钟同一个IP出现60次404
	warner404 := warner.NewWarner(60*time.Second, 60)
	warner404.ResetOnWarn = true
	warner404.On(func(ip string, _ warner.Count) {
		service.AlarmError("too many 404 request, client ip:" + ip)
	})
	// 定期清除warner中的过期数据
	go func() {
		for range time.NewTicker(5 * time.Minute).C {
			warner404.ClearExpired()
		}
	}()

	e.NotFoundHandler = func(resp http.ResponseWriter, req *http.Request) {
		ip := elton.GetRealIP(req)
		logger.Info("404",
			zap.String("ip", ip),
			zap.String("method", req.Method),
			zap.String("uri", req.RequestURI),
		)
		resp.Header().Set(elton.HeaderContentType, elton.MIMEApplicationJSON)
		resp.WriteHeader(http.StatusNotFound)
		_, _ = resp.Write([]byte(`{"statusCode": 404,"message": "Not found"}`))
		warner404.Inc(ip, 1)
	}

	// 捕捉panic异常，避免程序崩溃
	e.Use(eltonMid.NewRecover())

	e.Use(middleware.NewEntry())

	// 接口相关统计信息
	e.Use(eltonMid.NewStats(eltonMid.StatsConfig{
		OnStats: func(info *eltonMid.StatsInfo, c *elton.Context) {
			// ping 的日志忽略
			if info.URI == "/ping" {
				return
			}
			logger.Info("access log",
				zap.String("id", info.CID),
				zap.String("ip", info.IP),
				zap.String("sid", util.GetSessionID(c)),
				zap.String("method", info.Method),
				zap.String("uri", info.URI),
				zap.Int("status", info.Status),
				zap.String("consuming", info.Consuming.String()),
				zap.String("size", humanize.Bytes(uint64(info.Size))),
			)
		},
	}))

	// 错误处理，将错误转换为json响应
	e.Use(eltonMid.NewDefaultError())

	// IP限制
	e.Use(middleware.NewIPBlock())

	// 根据应用配置限制路由
	e.Use(middleware.NewRouterController())

	// 路由并发限制
	// routerLimitConfig := config.GetRouterConcurrentLimit()
	// if len(routerLimitConfig) != 0 {
	// 	e.Use(routerLimiter.New(routerLimiter.Config{
	// 		Limiter: routerLimiter.NewLocalLimiter(routerLimitConfig),
	// 	}))
	// }

	// etag与fresh的处理
	e.Use(eltonMid.NewDefaultFresh())
	e.Use(eltonMid.NewDefaultETag())

	// 对响应数据 c.Body 转换为相应的json响应
	e.Use(eltonMid.NewDefaultResponder())

	// 读取读取body的数的，转换为json bytes
	bodyparserConfig := eltonMid.BodyParserConfig{
		// 放宽5MB
		Limit: 5 * 1024 * 1024,
	}
	bodyparserConfig.AddDecoder(eltonMid.NewGzipDecoder())
	bodyparserConfig.AddDecoder(eltonMid.NewJSONDecoder())
	e.Use(eltonMid.NewBodyParser(bodyparserConfig))

	// 初始化路由
	for _, g := range router.GetGroups() {
		e.AddGroup(g)
	}

	err := dependServiceCheck()
	if err != nil {
		// 可以针对实际场景输出更多的日志信息
		logger.DPanic("exception",
			zap.Error(err),
		)
		panic(err)
	}
	logger.Info("start to linstening...",
		zap.String("listen", config.GetListen()),
	)
	// http1与http2均支持
	e.Server = &http.Server{
		Handler: h2c.NewHandler(e, &http2.Server{}),
	}

	err = e.ListenAndServe(config.GetListen())
	if err != nil {
		panic(err)
	}
}
