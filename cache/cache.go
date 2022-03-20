// Copyright 2021 tree xie
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

package cache

import (
	"time"

	"github.com/vicanso/tiny-site/config"
	"github.com/vicanso/tiny-site/helper"
	goCache "github.com/vicanso/go-cache"
	lruttl "github.com/vicanso/lru-ttl"
)

var redisCache = newRedisCache()
var redisCacheWithCompress = newCompressRedisCache()
var redisSession = newRedisSession()
var redisConfig = config.MustGetRedisConfig()

func newRedisCache() *goCache.RedisCache {
	opts := []goCache.RedisCacheOption{
		goCache.RedisCachePrefixOption(redisConfig.Prefix),
	}
	c := goCache.NewRedisCache(helper.RedisGetClient(), opts...)
	return c
}

func newCompressRedisCache() *goCache.RedisCache {
	// 大于10KB以上的数据压缩
	// 适用于数据量较大，而且数据内容重复较多的场景
	minCompressSize := 10 * 1024
	return goCache.NewLZ4RedisCache(
		helper.RedisGetClient(),
		minCompressSize,
		goCache.RedisCachePrefixOption(redisConfig.Prefix),
	)
}

func newRedisSession() *goCache.RedisSession {
	ss := goCache.NewRedisSession(helper.RedisGetClient())
	// 设置前缀
	ss.SetPrefix(redisConfig.Prefix + "ss:")
	return ss
}

// GetRedisCache get redis cache
func GetRedisCache() *goCache.RedisCache {
	return redisCache
}

// GetRedisCacheWithCompress get redis cache which will compress data
func GetRedisCacheWithCompress() *goCache.RedisCache {
	return redisCacheWithCompress
}

// GetRedisSession get redis session
func GetRedisSession() *goCache.RedisSession {
	return redisSession
}

// NewMultilevelCache create a new multilevel cache
func NewMultilevelCache(lruSize int, ttl time.Duration, prefix string) *lruttl.L2Cache {
	opts := []goCache.MultilevelCacheOption{
		goCache.MultilevelCacheRedisOption(redisCache),
		goCache.MultilevelCacheLRUSizeOption(lruSize),
		goCache.MultilevelCacheTTLOption(ttl),
		goCache.MultilevelCachePrefixOption(prefix),
	}
	return goCache.NewMultilevelCache(opts...)
}

// NewLRUCache new lru cache with ttl
func NewLRUCache(maxEntries int, defaultTTL time.Duration) *lruttl.Cache {
	return lruttl.New(maxEntries, defaultTTL)
}
