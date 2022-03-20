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

package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfigENV(t *testing.T) {
	assert := assert.New(t)
	originENV := env
	defer func() {
		env = originENV
	}()

	env = "test"
	assert.Equal(env, GetENV())
}

func TestBasicConfig(t *testing.T) {
	assert := assert.New(t)

	basicConfig := MustGetBasicConfig()
	assert.Equal("tiny-site", basicConfig.Name)
	assert.Equal(uint(1000), basicConfig.RequestLimit)
	assert.Equal(":7001", basicConfig.Listen)
}

func TestSessionConfig(t *testing.T) {
	assert := assert.New(t)

	sessionConfig := MustGetSessionConfig()
	assert.Equal(240*time.Hour, sessionConfig.TTL)
	assert.Equal("tiny-site", sessionConfig.Key)
	assert.Equal("/", sessionConfig.CookiePath)
	assert.Equal([]string{"cuttlefish", "secret"}, sessionConfig.Keys)
	assert.Equal("jt", sessionConfig.TrackKey)
}

func TestRedisConfig(t *testing.T) {
	assert := assert.New(t)

	redisConfig := MustGetRedisConfig()
	assert.Equal([]string{"127.0.0.1:6379"}, redisConfig.Addrs)
	assert.Equal("", redisConfig.Password)
	assert.Equal(200*time.Millisecond, redisConfig.Slow)
	assert.Equal(uint32(1000), redisConfig.MaxProcessing)
}

func TestMailConfig(t *testing.T) {
	assert := assert.New(t)

	mailConfig := MustGetMailConfig()
	assert.Equal("smtp.office365.com", mailConfig.Host)
	assert.Equal(587, mailConfig.Port)
	assert.Equal("tree.xie@outlook.com", mailConfig.User)
	assert.Equal("pass", mailConfig.Password)
}

func TestInfluxdbConfig(t *testing.T) {
	assert := assert.New(t)

	influxdbConfig := MustGetInfluxdbConfig()
	assert.Equal("http://127.0.0.1:8086", influxdbConfig.URI)
	assert.Equal("tiny-site", influxdbConfig.Bucket)
	assert.Equal("bigTree", influxdbConfig.Org)
	assert.NotEmpty(influxdbConfig.Token)
	assert.Equal(uint(100), influxdbConfig.BatchSize)
	assert.NotEmpty(influxdbConfig.FlushInterval)
	assert.False(influxdbConfig.Disabled)
}

func TestMustGetBasicConfig(t *testing.T) {
	assert := assert.New(t)

	databaseConfig := MustGetDatabaseConfig()
	assert.NotNil(databaseConfig.URI)
}

func TestMustGetLocationConfig(t *testing.T) {
	assert := assert.New(t)

	locationConfig := MustGetLocationConfig()
	assert.Equal("https://ip.npmtrend.com", locationConfig.BaseURL)
	assert.Equal(3*time.Second, locationConfig.Timeout)
}

func TestMustGetMinioConfig(t *testing.T) {
	assert := assert.New(t)

	minioConfig := MustGetMinioConfig()
	assert.Equal("127.0.0.1:9000", minioConfig.Endpoint)
	assert.Equal("origin", minioConfig.AccessKeyID)
	assert.Equal("test123456", minioConfig.SecretAccessKey)
	assert.False(minioConfig.SSL)
}
