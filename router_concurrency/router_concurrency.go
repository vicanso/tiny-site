package routerconcurrency

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/vicanso/elton"
	"github.com/vicanso/tiny-site/email"
	"github.com/vicanso/tiny-site/log"
	"go.uber.org/atomic"
	"go.uber.org/ratelimit"
)

type (
	// routerRateLimit 路由频率限制
	routerRateLimit struct {
		Limiter ratelimit.Limiter
	}
	// RouterConcurrency 路由并发配置
	RouterConcurrency struct {
		Router string `json:"router"`
		Max    uint32 `json:"max"`
		// 频率限制
		Rate int `json:"rate"`
		// 间隔
		Interval string `json:"interval"`

		// aotmic
		current       atomic.Uint32
		max           atomic.Uint32
		rateLimitDesc atomic.String
		// limit 保存routerRateLimit对象
		limit atomic.Value
	}
	// rcLimiter 路由请求限制
	rcLimiter struct {
		m map[string]*RouterConcurrency
	}
)

var (
	currentRCLimiter = &rcLimiter{}

	// 无频率限制
	routerRateUnlimited = &routerRateLimit{
		Limiter: ratelimit.NewUnlimited(),
	}
)

// SetRateLimiter 设置频率限制
func (rc *RouterConcurrency) update(item *RouterConcurrency) {
	// 设置并发请求量
	rc.max.Store(item.Max)
	rate := item.Rate
	interval := item.Interval
	// 获取rate limit配置，如果有调整则需要重新设置
	rateDesc := fmt.Sprintf("%d-%s", rate, interval)
	if rateDesc == rc.rateLimitDesc.Load() {
		return
	}
	d, _ := time.ParseDuration(interval)
	rc.rateLimitDesc.Store(rateDesc)
	// 如果未设置限制，则使用无限制频率
	// 如果未设置时长
	if rate <= 0 || d == 0 {
		rc.limit.Store(routerRateUnlimited)
		return
	}
	rrl := &routerRateLimit{
		Limiter: ratelimit.New(rate, ratelimit.Per(d)),
	}
	rc.limit.Store(rrl)
}

// Take 执行一次频率限制，此执行会根据当时频率延时
func (rc *RouterConcurrency) Take() {
	p := rc.limit.Load()
	if p == nil {
		return
	}
	limit, _ := p.(*routerRateLimit)
	if limit == nil {
		return
	}
	limit.Limiter.Take()
}

// IncConcurrency 当前路由处理数+1
func (l *rcLimiter) IncConcurrency(key string) (uint32, uint32) {
	// 该map仅初始化一次，因此无需要考虑锁
	r, ok := l.m[key]
	if !ok {
		return 0, 0
	}
	current := r.current.Inc()
	max := r.max.Load()
	// 如果设置为0或已超出最大并发限制，则直接返回
	if max == 0 || current > max {
		return current, max
	}
	r.Take()
	return current, max
}

// DecConcurrency 当前路由处理数-1
func (l *rcLimiter) DecConcurrency(key string) {
	r, ok := l.m[key]
	if !ok {
		return
	}
	r.current.Dec()
}

// GetConcurrency 获取当前路由处理数
func (l *rcLimiter) GetConcurrency(key string) uint32 {
	r, ok := l.m[key]
	if !ok {
		return 0
	}
	return r.current.Load()
}

// GetStats 获取统计
func (l *rcLimiter) GetStats() map[string]uint32 {
	result := make(map[string]uint32)
	for key, r := range l.m {
		result[key] = r.current.Load()
	}
	return result
}

// InitLimiter 初始路由并发限制
func InitLimiter(routers []elton.RouterInfo) {
	m := make(map[string]*RouterConcurrency)
	for _, item := range routers {
		m[item.Method+" "+item.Route] = &RouterConcurrency{}
	}
	currentRCLimiter.m = m
}

// GetLimiter 获取路由并发限制器
func GetLimiter() *rcLimiter {
	return currentRCLimiter
}

// Update 更新路由并发数
func Update(arr []string) {
	concurrencyConfigList := make([]*RouterConcurrency, 0)
	for _, str := range arr {
		v := &RouterConcurrency{}
		err := json.Unmarshal([]byte(str), v)
		if err != nil {
			log.Error(context.Background()).
				Err(err).
				Msg("router concurrency config is invalid")
			email.AlarmError(context.Background(), "router concurrency config is invalid:"+err.Error())
			continue
		}
		concurrencyConfigList = append(concurrencyConfigList, v)
	}
	for key, r := range currentRCLimiter.m {
		found := false
		for _, item := range concurrencyConfigList {
			if item.Router == key {
				found = true
				r.update(item)
			}
		}
		// 如果未配置，则设置为限制0（无限制）
		if !found {
			r.max.Store(0)
		}
	}
}

// List 获取路由并发限制数
func List() map[string]uint32 {
	result := make(map[string]uint32)
	for key, r := range currentRCLimiter.m {
		v := r.max.Load()
		if v != 0 {
			result[key] = v
		}
	}
	return result
}
