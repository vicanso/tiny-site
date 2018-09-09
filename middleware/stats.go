package middleware

import (
	"net/url"
	"sync/atomic"
	"time"

	"github.com/kataras/iris"
)

type (
	// OnStats on stats function
	OnStats func(*StatsInfo)
	// StatsConfig stats的配置
	StatsConfig struct {
		OnStats OnStats
	}
	// StatsInfo 统计信息
	StatsInfo struct {
		IP         string
		TrackID    string
		Account    string
		Method     string
		Route      string
		URI        string
		StatusCode int
		Consuming  int
		Type       int
		Connecting uint32
	}
)

// NewStats 请求统计
func NewStats(conf StatsConfig) iris.Handler {
	var connectingCount uint32
	return func(ctx iris.Context) {
		atomic.AddUint32(&connectingCount, 1)
		startedAt := time.Now().UnixNano()
		req := ctx.Request()
		uri, _ := url.QueryUnescape(req.RequestURI)
		if uri == "" {
			uri = req.RequestURI
		}
		ip := ctx.RemoteAddr()
		trackID := getTrackID(ctx)
		route := ctx.GetCurrentRoute()

		ctx.Next()
		consuming := int(time.Now().UnixNano()-startedAt) / int(time.Millisecond)
		statusCode := ctx.GetStatusCode()

		info := &StatsInfo{
			URI:        uri,
			StatusCode: statusCode,
			Consuming:  consuming,
			Type:       statusCode / 100,
			Connecting: connectingCount,
			IP:         ip,
			TrackID:    trackID,
			Account:    getAccount(ctx),
		}
		if route != nil {
			info.Method = route.Method()
			info.Route = route.Path()
		}
		if conf.OnStats != nil {
			conf.OnStats(info)
		}
		atomic.AddUint32(&connectingCount, ^uint32(0))
	}
}
