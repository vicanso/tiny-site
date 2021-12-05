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
	"fmt"
	"strconv"
	"time"

	"github.com/vicanso/elton"
	"github.com/vicanso/elton/middleware"
	"github.com/vicanso/tiny-site/cache"
	"github.com/vicanso/tiny-site/log"
	"github.com/vicanso/hes"
)

const (
	concurrentLimitKeyPrefix = "midConcurrentLimit"
	ipLimitKeyPrefix         = "midIPLimit"
	errorLimitKeyPrefix      = "midErrorLimit"
	errLimitCategory         = "requestLimit"
)

var redisSrv = cache.GetRedisCache()

type (
	// KeyGenerator key generator
	KeyGenerator func(c *elton.Context) string
)

// createConcurrentLimitLock 创建并发限制的lock函数
func createConcurrentLimitLock(prefix string, ttl time.Duration, withDone bool) middleware.ConcurrentLimiterLock {
	return func(key string, c *elton.Context) (bool, func(), error) {
		ctx := c.Context()
		k := concurrentLimitKeyPrefix + "-" + prefix + "-" + key
		var done func()
		if withDone {
			success, redisDone, err := redisSrv.LockWithDone(ctx, k, ttl)
			done = func() {
				err := redisDone()
				if err != nil {
					log.Error(ctx).
						Str("category", "redisDelFail").
						Str("key", k).
						Err(err).
						Msg("")
				}
			}
			return success, done, err
		}
		success, err := redisSrv.Lock(ctx, k, ttl)
		if err != nil {
			return false, done, err
		}
		return success, done, nil
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

// NewConcurrentLimitWithDone 创建并发限制中间件，完成时则执行done清除
func NewConcurrentLimitWithDone(keys []string, ttl time.Duration, prefix string) elton.Handler {
	return middleware.NewConcurrentLimiter(middleware.ConcurrentLimiterConfig{
		NotAllowEmpty: true,
		Lock:          createConcurrentLimitLock(prefix, ttl, true),
		Keys:          keys,
	})
}

// NewIPLimit 创建IP限制中间件
func NewIPLimit(maxCount int64, ttl time.Duration, prefix string) elton.Handler {
	return func(c *elton.Context) error {
		ctx := c.Context()
		key := ipLimitKeyPrefix + "-" + prefix + "-" + c.RealIP()
		count, err := redisSrv.IncWith(ctx, key, 1, ttl)
		if err != nil {
			return err
		}
		if count > maxCount {
			return hes.New(fmt.Sprintf("请求过于频繁，请稍候再试！(%d/%d)", count, maxCount), errLimitCategory)
		}
		return c.Next()
	}
}

// NewErrorLimit 创建出错限制中间件
func NewErrorLimit(maxCount int64, ttl time.Duration, fn KeyGenerator) elton.Handler {
	return func(c *elton.Context) error {
		ctx := c.Context()
		key := errorLimitKeyPrefix + "-" + fn(c)
		result, err := redisSrv.GetIgnoreNilErr(ctx, key)
		if err != nil {
			return err
		}
		count, _ := strconv.Atoi(string(result))
		// 因为count是处理完才inc，因此增加等于的判断
		if int64(count) >= maxCount {
			return hes.New(fmt.Sprintf("请求过于频繁，请稍候再试！(%d/%d)", count, maxCount), errLimitCategory)
		}
		err = c.Next()
		// 如果出错，则出错次数+1
		if err != nil {
			_, _ = redisSrv.IncWith(ctx, key, 1, ttl)
		}
		return err
	}
}
