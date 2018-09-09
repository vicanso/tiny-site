package middleware

import (
	"sync/atomic"
	"time"

	"github.com/kataras/iris"
	"github.com/vicanso/tiny-site/global"
	"github.com/vicanso/tiny-site/util"
)

type (
	// LimiterConfig limiter配置
	LimiterConfig struct {
		Max uint32
	}
)

func resetApplicationStatus(d time.Duration) {
	// 等待后将程序重置为可用
	ticker := time.NewTicker(d)
	go func() {
		select {
		case <-ticker.C:
			// TODO 是否需要记录相关resume记录
			logger := util.GetLogger()
			logger.Info("application resume")
			global.StartApplication()
		}
	}()
}

// NewLimiter 连接限制中间件
func NewLimiter(conf LimiterConfig) iris.Handler {
	var count uint32
	return func(ctx iris.Context) {
		defer func() {
			atomic.AddUint32(&count, ^uint32(0))
		}()
		v := atomic.AddUint32(&count, 1)
		if v > conf.Max {
			resErr(ctx, util.ErrTooManyRequest)
			// 如果多并发，还是会导致多个reset，影响不大，忽略
			if global.IsApplicationRunning() {
				global.PauseApplication()
				resetApplicationStatus(time.Second * 10)
			}
			return
		}
		ctx.Next()
	}
}
