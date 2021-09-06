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
	"context"
	"time"

	"github.com/go-redis/redis"
	"github.com/vicanso/elton"
	session "github.com/vicanso/elton-session"
	"github.com/vicanso/tiny-site/config"
	"github.com/vicanso/tiny-site/helper"
	"github.com/vicanso/tiny-site/util"
)

type (
	// RedisStore redis store for session
	RedisStore struct {
		client *redis.Client
		prefix string
	}
)

func (rs *RedisStore) getKey(key string) string {
	return rs.prefix + key
}

// Get get the session from redis
func (rs *RedisStore) Get(_ context.Context, key string) ([]byte, error) {
	buf, err := rs.client.Get(rs.getKey(key)).Bytes()
	if err == redis.Nil {
		return buf, nil
	}
	return buf, err
}

// Set set the session to redis
func (rs *RedisStore) Set(_ context.Context, key string, data []byte, ttl time.Duration) error {
	return rs.client.Set(rs.getKey(key), data, ttl).Err()
}

// Destroy remove the session from redis
func (rs *RedisStore) Destroy(_ context.Context, key string) error {
	return rs.client.Del(rs.getKey(key)).Err()
}

// NewRedisStore create new redis store instance
func NewRedisStore(client *redis.Client, prefix string) *RedisStore {
	rs := &RedisStore{}
	rs.client = client
	rs.prefix = prefix
	return rs
}

// NewSession new session middleware
func NewSession() elton.Handler {
	client := helper.RedisGetClient()
	if client == nil {
		panic("session store need redis client")
	}
	store := NewRedisStore(client, "ss-")
	scf := config.GetSessionConfig()
	return session.NewByCookie(session.CookieConfig{
		Store:   store,
		Signed:  true,
		Expired: scf.TTL,
		GenID: func() string {
			return util.GenUlid()
		},
		Name:     scf.Key,
		Path:     scf.CookiePath,
		MaxAge:   int(scf.TTL.Seconds()),
		HttpOnly: true,
	})
}
