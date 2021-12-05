// Copyright 2020 tree xie
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

/*
Package main tiny-site

	接口出错统一使用如下格式：{"category": "出错类别", "message": "出错信息", "code": "出错代码", "exception": true}，
	其中category与code字段为可选，当处理出错时，HTTP的响应状态码为`4xx`与`5xx`。
	如果exception为true则表示此出错未预期出错，一般需要修复。
	其中`4xx`表示客户端参数等异常出错，而`5xx`则表示服务处理异常。


	常见出错类别：
	`validate`：表示参数校验失败，接口传参不符合约束条件

	常用pattern:

	`xLimit`:  >=1 && <=100 的整数数

Host: 127.0.0.1:7001
Version: 1.0.0
Schemes: http

Consumes:
- application/json

Produces:
- application/json

swagger:meta
*/
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"sync"
	"syscall"
	"time"

	warner "github.com/vicanso/count-warner"
	"github.com/vicanso/elton"
	compress "github.com/vicanso/elton-compress"
	M "github.com/vicanso/elton/middleware"
	"github.com/vicanso/tiny-site/config"
	_ "github.com/vicanso/tiny-site/controller"
	"github.com/vicanso/tiny-site/cs"
	"github.com/vicanso/tiny-site/email"
	"github.com/vicanso/tiny-site/helper"
	"github.com/vicanso/tiny-site/log"
	"github.com/vicanso/tiny-site/middleware"
	"github.com/vicanso/tiny-site/profiler"
	"github.com/vicanso/tiny-site/router"
	routerconcurrency "github.com/vicanso/tiny-site/router_concurrency"
	routermock "github.com/vicanso/tiny-site/router_mock"
	_ "github.com/vicanso/tiny-site/schedule"
	"github.com/vicanso/tiny-site/service"
	"github.com/vicanso/tiny-site/util"
	"github.com/vicanso/hes"
	"go.uber.org/automaxprocs/maxprocs"
)

var (
	// Version 应用版本号
	Version string
	// BuildAt 构建时间
	BuildedAt string
)

// closeDepends 程序关闭时关闭依赖的服务
var closeDepends func()

func init() {
	// 启动出错中的caller记录
	hes.EnableCaller(true)

	// 替换出错信息中的file中的目录
	basicConfig := config.MustGetBasicConfig()
	reg := regexp.MustCompile(fmt.Sprintf(`\S*/%s/`, basicConfig.Name))
	hes.SetFileConvertor(func(file string) string {
		return reg.ReplaceAllString(file, "")
	})

	_, _ = maxprocs.Set(maxprocs.Logger(func(format string, args ...interface{}) {
		value := fmt.Sprintf(format, args...)
		log.Info(context.Background()).
			Msg(value)
	}))
	service.SetApplicationVersion(Version)
	service.SetApplicationBuildedAt(BuildedAt)
	closeOnce := sync.Once{}
	closeDepends = func() {
		closeOnce.Do(func() {
			// 关闭influxdb，flush统计数据
			helper.GetInfluxDB().Close()
			_ = helper.EntGetClient().Close()
			_ = helper.RedisGetClient().Close()
		})
	}
}

// 是否用户主动关闭
var closedByUser = false

func gracefulClose(e *elton.Elton) {
	log.Info(context.Background()).Msg("start to graceful close")
	// 设置状态为退出中，/ping请求返回出错，反向代理不再转发流量
	service.SetApplicationStatus(service.ApplicationStatusStopping)
	// docker 在10秒内退出，因此设置5秒后退出
	time.Sleep(5 * time.Second)
	// 所有新的请求均返回出错
	e.GracefulClose(3 * time.Second)
	closeDepends()
	os.Exit(0)
}

// watchForClose 监听信号关闭程序
func watchForClose(e *elton.Elton) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			log.Info(context.Background()).
				Str("signal", s.String()).
				Msg("server will be closed")
			closedByUser = true
			gracefulClose(e)
		}
	}()
}

// devWaitForExit 开发环境退出
func devWaitForExit(e *elton.Elton) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	go func() {
		for range c {
			closeDepends()
			e.Close()
			os.Exit(1)
		}
	}()
}

// 相关依赖服务的校验，主要是数据库等
func dependServiceCheck() (err error) {
	err = helper.RedisPing()
	if err != nil {
		return
	}
	err = helper.EntPing()
	if err != nil {
		return
	}
	// 程序启动后再执行init schema
	err = helper.EntInitSchema()
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

func newOnErrorHandler(e *elton.Elton) {

	// 未处理的error才会触发
	// 如果1分钟出现超过5次未处理异常
	// exception的warner只有一个key，因此无需定时清除
	exceptionWarner := warner.NewWarner(5*time.Minute, 5)
	exceptionWarner.On(func(_ string, _ int) {
		email.AlarmError(context.Background(), "too many uncaught exception")
	})
	// 只有未被处理的error才会触发此回调
	e.OnError(func(c *elton.Context, err error) {
		he := hes.Wrap(err)

		stack := util.GetStack(5)
		if len(stack) != 0 {
			stack = stack[1:]
		}
		ip := c.RealIP()
		uri := c.Request.RequestURI
		he.AddExtra("route", c.Route)
		he.AddExtra("uri", uri)

		// 记录exception
		service.GetInfluxSrv().Write(cs.MeasurementException, map[string]string{
			cs.TagCategory: "routeError",
			cs.TagRoute:    c.Route,
		}, map[string]interface{}{
			cs.FieldIP:  ip,
			cs.FieldURI: uri,
		})

		// 可以针对实际场景输出更多的日志信息
		log.Error(c.Context()).
			Str("category", "exception").
			Str("ip", ip).
			Str("route", c.Route).
			Str("uri", uri).
			Strs("stack", stack).
			Msg("")

		exceptionWarner.Inc("exception", 1)
		// panic类的异常都graceful close
		if he.Category == M.ErrRecoverCategory {

			email.AlarmError(c.Context(), "panic recover:"+string(he.ToJSON()))
			// 由于此处的error由请求触发的，因为要另外启动一个goroutine重启，避免影响当前处理
			go gracefulClose(e)
		}
	})
}

func main() {
	profiler.MustStartPyroscope()
	e := elton.New()
	// 记录server中连接的状态变化
	e.Server.ConnState = service.GetHTTPServerConnState()
	e.Server.ErrorLog = log.NewHTTPServerLogger()

	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
			log.Error(context.Background()).
				Str("category", "panic").
				Err(err).
				Msg("")
			email.AlarmError(context.Background(), fmt.Sprintf("panic recover:%v", r))
			// panic类的异常都graceful close
			gracefulClose(e)
		}
	}()

	basicConfig := config.MustGetBasicConfig()
	defer closeDepends()
	// 非开发环境，监听信号退出
	if !util.IsDevelopment() {
		watchForClose(e)
		go func() {
			pidData := []byte(strconv.Itoa(os.Getpid()))
			err := ioutil.WriteFile(basicConfig.PidFile, pidData, 0600)
			if err != nil {
				log.Error(context.Background()).
					Err(err).
					Msg("write pid fail")
			}
		}()
	} else {
		devWaitForExit(e)
	}

	newOnErrorHandler(e)
	// 启用耗时跟踪
	if util.IsDevelopment() {
		e.EnableTrace = true
	}
	e.SignedKeys = service.GetSignedKeys()
	e.OnTrace(func(c *elton.Context, infos elton.TraceInfos) {
		// 设置server timing
		c.ServerTiming(infos, "tiny-site-")
	})
	// 若需要唯一值，可使用ulid或uuid
	e.GenerateID = func() string {
		return util.RandomString(6)
	}

	// 自定义404与405的处理
	e.NotFoundHandler = middleware.NewNotFoundHandler()
	e.MethodNotAllowedHandler = middleware.NewMethodNotAllowedHandler()

	// 前缀处理
	if len(basicConfig.Prefixes) != 0 {
		e.Pre(middleware.NewPrefixHandler(basicConfig.Prefixes))
	}

	// 捕捉panic异常，避免程序崩溃
	e.UseWithName(M.NewRecover(), "recover")

	// 入口设置
	e.UseWithName(middleware.NewEntry(service.IncreaseConcurrency, service.DecreaseConcurrency), "entry")

	// 接口相关统计信息
	e.UseWithName(middleware.NewStats(), "stats")

	// 出错转换为json（出错处理应该在stats之后，这样stats中才可获取到正确的http status code)
	e.UseWithName(middleware.NewError(), "error")

	// 仅将timeout设置给context，后续调用如果无依赖于context
	// 则不会超时
	// 后续再考虑是否增加select
	e.UseWithName(func(c *elton.Context) error {
		ctx, cancel := context.WithTimeout(c.Context(), basicConfig.Timeout)
		defer cancel()
		c.WithContext(ctx)
		return c.Next()
	}, "timeout")

	// 限制最大请求量
	if basicConfig.RequestLimit != 0 {
		e.UseWithName(M.NewGlobalConcurrentLimiter(M.GlobalConcurrentLimiterConfig{
			Max: uint32(basicConfig.RequestLimit),
		}), "requestLimit")
	}
	// tracer中间件在最大请求限制中间件之后，保证进入tracer的goroutine不要过多
	// e.UseWithName(tracer.New(), "tracer")

	// 配置只针对snappy与zstd压缩（主要用于减少内网线路带宽，对外的压缩由前置反向代理完成）
	compressMinLength := 2 * 1024
	compressConfig := M.NewCompressConfig(
		&compress.SnappyCompressor{
			MinLength: compressMinLength,
		},
		&compress.ZstdCompressor{
			MinLength: compressMinLength,
		},
	)
	e.UseWithName(M.NewCompress(compressConfig), "compress")

	// IP限制
	e.UseWithName(middleware.NewIPBlocker(service.IsBlockIP), "ipBlocker")

	// 根据配置对路由mock返回
	e.UseWithName(middleware.NewRouterMocker(routermock.Get), "routerMocker")

	// 路由并发限制
	e.UseWithName(M.NewRCL(M.RCLConfig{
		Limiter: routerconcurrency.GetLimiter(),
	}), "rcl")

	// eTag与fresh的处理
	e.UseWithName(M.NewDefaultFresh(), "fresh").
		UseWithName(M.NewDefaultETag(), "eTag")

	// 读取读取body的数的，转换为json bytes
	e.UseWithName(M.NewDefaultBodyParser(), "bodyParser")

	// 拦截
	e.UseWithName(middleware.NewInterceptor(), "interceptor")

	// 对响应数据 c.Body 转换为相应的json响应
	e.UseWithName(M.NewDefaultResponder(), "responder")

	// 初始化路由
	e.AddGroup(router.GetGroups()...)

	// 初始化路由并发限制配置
	routerconcurrency.InitLimiter(e.GetRouters())

	err := dependServiceCheck()
	if err != nil {
		email.AlarmError(context.Background(), "check depend service fail, "+err.Error())
		log.Error(context.Background()).
			Str("category", "depFail").
			Err(err).
			Msg("")
		return
	}

	service.SetApplicationStatus(service.ApplicationStatusRunning)

	// http1与http2均支持
	// 一般后端服务可以不需要启用
	// e.Server = &http.Server{
	// 	Handler: h2c.NewHandler(e, &http2.Server{}),
	// }
	log.Info(context.Background()).Msg("server will listen on " + basicConfig.Listen)
	err = e.ListenAndServe(basicConfig.Listen)
	// 如果出错而且非主动关闭，则发送告警
	if err != nil && !closedByUser {
		email.AlarmError(context.Background(), "listen and serve fail, "+err.Error())
		log.Error(context.Background()).
			Err(err).
			Msg("")
	}
}
