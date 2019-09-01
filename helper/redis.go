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

package helper

import (
	"github.com/go-redis/redis"
	"go.uber.org/zap"

	"github.com/vicanso/tiny-site/config"
)

var (
	redisClient *redis.Client
)

func init() {
	options, err := config.GetRedisConfig()
	if err != nil {
		panic(err)
	}
	logger.Info("connect to redis",
		zap.String("addr", options.Addr),
		zap.Int("db", options.DB),
	)
	redisClient = redis.NewClient(&redis.Options{
		Addr:     options.Addr,
		Password: options.Password,
		DB:       options.DB,
	})
}

// RedisGetClient get redis client
func RedisGetClient() *redis.Client {
	return redisClient
}
