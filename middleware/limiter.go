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

package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/vicanso/elton"
	"github.com/vicanso/elton/middleware"
	"github.com/vicanso/hes"
	"github.com/vicanso/tiny-site/log"
	"github.com/vicanso/tiny-site/service"
	"go.uber.org/zap"
)

type (
	// KeyGenerator key generator
	KeyGenerator func(c *elton.Context) string
)

const (
	concurrentLimitKeyPrefix = "midConcurrentLimit"
	ipLimitKeyPrefix         = "midIPLimit"
	errorLimitKeyPrefix      = "midErrorLimit"
	errLimitCategory         = "requestLimit"
)

var (
	errTooFrequently = &hes.Error{
		StatusCode: http.StatusBadRequest,
		Message:    "请求过于频繁，请稍候再试！",
		Category:   errLimitCategory,
	}
	redisSrv = new(service.RedisSrv)

	logger = log.Default()
)

// createConcurrentLimitLock 创建并发限制的lock函数
func createConcurrentLimitLock(prefix string, ttl time.Duration, withDone bool) middleware.ConcurrentLimiterLock {
	return func(key string, _ *elton.Context) (success bool, done func(), err error) {
		k := concurrentLimitKeyPrefix + "-" + prefix + "-" + key
		done = nil
		if withDone {
			success, redisDone, err := redisSrv.LockWithDone(k, ttl)
			done = func() {
				err := redisDone()
				if err != nil {
					logger.Error("redis done fail",
						zap.String("key", k),
						zap.Error(err),
					)
				}
			}
			return success, done, err
		}
		success, err = redisSrv.Lock(k, ttl)
		return
	}
}

// NewConcurrentLimit 创建并发限制的中间件
func NewConcurrentLimit(keys []string, ttl time.Duration, prefix string) elton.Handler {
	return middleware.NewConcurrentLimiter(middleware.ConcurrentLimiterConfig{
		NotAllowEmpty: true,
		Lock:          createConcurrentLimitLock(prefix, ttl, false),
		Keys:          keys,
	})
}

// NewConcurrentLimitWithDone 创建并发限制中间件，且带done函数
func NewConcurrentLimitWithDone(keys []string, ttl time.Duration, prefix string) elton.Handler {
	return middleware.NewConcurrentLimiter(middleware.ConcurrentLimiterConfig{
		NotAllowEmpty: true,
		Lock:          createConcurrentLimitLock(prefix, ttl, true),
		Keys:          keys,
	})
}

// NewIPLimit 创建IP限制中间件
func NewIPLimit(maxCount int64, ttl time.Duration, prefix string) elton.Handler {
	return func(c *elton.Context) (err error) {
		key := ipLimitKeyPrefix + "-" + prefix + "-" + c.RealIP()
		count, err := redisSrv.IncWithTTL(key, ttl)
		if err != nil {
			return
		}
		if count > maxCount {
			err = errTooFrequently
			return
		}
		return c.Next()
	}
}

// NewErrorLimit 创建出错限制中间件
func NewErrorLimit(maxCount int64, ttl time.Duration, fn KeyGenerator) elton.Handler {
	return func(c *elton.Context) (err error) {
		key := errorLimitKeyPrefix + "-" + fn(c)
		result, err := redisSrv.GetIgnoreNilErr(key)
		if err != nil {
			return
		}
		count, _ := strconv.Atoi(result)
		// 因为count是处理完才inc，因此增加等于的判断
		if int64(count) >= maxCount {
			err = errTooFrequently
			return
		}
		err = c.Next()
		// 如果出错，则出错次数+1
		if err != nil {
			_, _ = redisSrv.IncWithTTL(key, ttl)
		}
		return
	}
}
