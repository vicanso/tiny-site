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
	"github.com/vicanso/tiny-site/log"
	"github.com/vicanso/tiny-site/service"
	"github.com/vicanso/hes"
	"go.uber.org/zap"

	concurrentLimiter "github.com/vicanso/elton-concurrent-limiter"
)

const (
	concurrentLimitKeyPrefix = "mid-concurrent-limit"
	ipLimitKeyPrefix         = "mid-ip-limit"
	errorLimitKeyPrefix      = "mid-error-limit"
	errLimitCategory         = "request-limit"
)

var (
	errTooFrequently = &hes.Error{
		StatusCode: http.StatusBadRequest,
		Message:    "request to frequently",
		Category:   errLimitCategory,
	}
	redisSrv = new(service.RedisSrv)
)

type (
	// KeyGenerator key generator
	KeyGenerator func(c *elton.Context) string
)

// createConcurrentLimitLock 创建并发限制的lock函数
func createConcurrentLimitLock(prefix string, ttl time.Duration, withDone bool) concurrentLimiter.Lock {
	return func(key string, _ *elton.Context) (success bool, done func(), err error) {
		k := concurrentLimitKeyPrefix + "-" + prefix + "-" + key
		done = nil
		if withDone {
			success, redisDone, err := redisSrv.LockWithDone(k, ttl)
			done = func() {
				err := redisDone()
				if err != nil {
					log.Default().Error("redis done fail",
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

// NewConcurrentLimit create a concurrent limit
func NewConcurrentLimit(keys []string, ttl time.Duration, prefix string) elton.Handler {
	return concurrentLimiter.New(concurrentLimiter.Config{
		Lock: createConcurrentLimitLock(prefix, ttl, false),
		Keys: keys,
	})
}

// NewIPLimit create a limit middleware by ip address
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

// NewErrorLimit create a error limit middleware
func NewErrorLimit(maxCount int64, ttl time.Duration, fn KeyGenerator) elton.Handler {
	return func(c *elton.Context) (err error) {
		key := errorLimitKeyPrefix + "-" + fn(c)
		result, err := redisSrv.Get(key)
		if err != nil {
			return
		}
		count, _ := strconv.Atoi(result)
		if int64(count) > maxCount {
			err = errTooFrequently
			return
		}
		err = c.Next()
		// 如果出错，则出错次数+1
		if err != nil {
			redisSrv.IncWithTTL(key, ttl)
		}
		return
	}
}
