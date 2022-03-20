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

// 公共的处理函数，包括程序基本信息、性能指标等

package controller

import (
	"bytes"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/vicanso/elton"
	"github.com/vicanso/tiny-site/asset"
	"github.com/vicanso/tiny-site/config"
	"github.com/vicanso/tiny-site/profiler"
	"github.com/vicanso/tiny-site/request"
	"github.com/vicanso/tiny-site/router"
	"github.com/vicanso/tiny-site/schema"
	"github.com/vicanso/tiny-site/service"
	"github.com/vicanso/tiny-site/util"
	"github.com/vicanso/hes"
)

type commonCtrl struct{}

// 响应相关定义
type (
	// applicationInfoResp 应用信息响应
	applicationInfoResp struct {
		// 版本号
		Version string `json:"version"`
		// 构建时间
		BuildedAt string `json:"buildedAt"`
		// 运行时长
		Uptime string `json:"uptime"`
		// os类型
		OS string `json:"os"`
		// go版本
		GO string `json:"go"`
		// 架构类型
		ARCH string `json:"arch"`
		// 运行环境配置
		ENV string `json:"env"`
	}
	// routersResp 路由列表响应
	routersResp struct {
		// 路由信息
		Routers []elton.RouterInfo `json:"routers"`
	}
	// statusListResp 状态列表响应
	statusListResp struct {
		Statuses []*schema.StatusInfo `json:"statuses"`
	}
	// randomKeysResp 随机字符
	randomKeysResp struct {
		Keys []string `json:"keys"`
	}
	// httpStatsListResp http性能统计响应
	httpStatsListResp struct {
		StatsList []*request.InstanceStats `json:"statsList"`
	}
)

const (
	errCommonCategory = "common"
)

var (
	// applicationStartedAt 应用启动时间
	applicationStartedAt = time.Now()
)

func init() {
	ctrl := commonCtrl{}
	router.NewGroup("").GET("/ping", ctrl.ping)
	g := router.NewGroup("/commons")

	g.GET("/application", ctrl.getApplicationInfo)
	g.GET("/routers", ctrl.getRouters)
	g.GET("/captcha", ctrl.getCaptcha)
	g.GET("/performance", ctrl.getPerformance)
	g.GET("/schema-statuses", ctrl.listStatus)
	g.GET("/random-keys", ctrl.getRandomKeys)
	// 获取系统prof指标
	g.GET(
		"/prof",
		loadUserSession,
		shouldBeAdmin,
		ctrl.getProf,
	)
	// 获取接口文档
	g.GET(
		"/api",
		ctrl.getAPI,
	)
	// 获取http实例性能指标
	g.GET(
		"/http-stats",
		ctrl.listHTTPInstanceStats,
	)
}

// ping 用于检测服务是否可用
func (*commonCtrl) ping(c *elton.Context) error {
	if !service.ApplicationIsRunning() {
		return hes.NewWithStatusCode("应用服务不可用", http.StatusServiceUnavailable, errCommonCategory)
	}
	c.BodyBuffer = bytes.NewBufferString("pong")
	return nil
}

// getApplicationInfo 获取应用信息
func (*commonCtrl) getApplicationInfo(c *elton.Context) error {
	c.CacheMaxAge(time.Minute)
	c.Body = &applicationInfoResp{
		Version:   service.GetApplicationVersion(),
		BuildedAt: service.GetApplicationBuildedAt(),
		Uptime:    humanize.Time(applicationStartedAt),
		OS:        runtime.GOOS,
		GO:        runtime.Version(),
		ARCH:      runtime.GOARCH,
		ENV:       config.GetENV(),
	}
	return nil
}

// getRouters 获取系统的路由
func (*commonCtrl) getRouters(c *elton.Context) error {
	c.CacheMaxAge(time.Minute)
	c.Body = &routersResp{
		Routers: c.Elton().GetRouters(),
	}
	return nil
}

// getCaptcha 获取图形验证码
func (*commonCtrl) getCaptcha(c *elton.Context) error {
	bgColor := c.QueryParam("bg")
	fontColor := c.QueryParam("color")
	if bgColor == "" {
		bgColor = "255,255,255"
	}
	if fontColor == "" {
		fontColor = "102,102,102"
	}
	info, err := service.GetCaptcha(c.Context(), fontColor, bgColor)
	if err != nil {
		return err
	}
	// 防止此字段未设置好，序列化至前端
	info.Value = ""
	c.NoStore()
	c.Body = info
	return nil
}

// getPerformance 获取应用性能指标
func (*commonCtrl) getPerformance(c *elton.Context) error {
	p := service.GetPerformance(c.Context())
	c.Body = p
	return nil
}

// listStatus 获取状态列表
func (*commonCtrl) listStatus(c *elton.Context) error {
	c.CacheMaxAge(5 * time.Minute)
	c.Body = &statusListResp{
		Statuses: schema.GetStatusList(),
	}
	return nil
}

// getRandomKeys 获取随机字符串
func (*commonCtrl) getRandomKeys(c *elton.Context) error {
	n, _ := strconv.Atoi(c.QueryParam("n"))
	size, _ := strconv.Atoi(c.QueryParam("size"))
	if size < 1 {
		size = 10
	}
	if n < 1 {
		n = 1
	}
	result := make([]string, n)
	for index := 0; index < n; index++ {
		result[index] = util.RandomString(size)
	}
	c.Body = &randomKeysResp{
		Keys: result,
	}
	return nil
}

// getProf 获取prof信息
func (*commonCtrl) getProf(c *elton.Context) error {
	d := 30 * time.Second
	v := c.QueryParam("d")
	if v != "" {
		d1, err := time.ParseDuration(v)
		if err != nil {
			return err
		}
		d = d1
	}
	result, err := profiler.GetProf(d)
	if err != nil {
		return err
	}
	c.SetHeader(elton.HeaderContentType, elton.MIMEBinary)
	c.SetHeader("Content-Disposition", `attachment; filename="gprof"`)
	c.BodyBuffer = result
	return nil
}

// getAPI 获取API信息
func (*commonCtrl) getAPI(c *elton.Context) error {
	file := "api.yml"
	buf, err := asset.GetFS().ReadFile(file)
	if err != nil {
		return err
	}
	c.SetHeader(elton.HeaderContentType, "text/vnd.yaml;charset=utf-8")

	c.BodyBuffer = bytes.NewBuffer(buf)
	return nil
}

// listHTTPInstanceStats 获取http实例的性能统计
func (*commonCtrl) listHTTPInstanceStats(c *elton.Context) error {
	stats := request.GetHTTPStats()
	c.CacheMaxAge(5 * time.Minute)
	c.Body = &httpStatsListResp{
		StatsList: stats,
	}
	return nil
}
