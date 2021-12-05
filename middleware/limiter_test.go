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
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
	"github.com/vicanso/tiny-site/util"
	"github.com/vicanso/hes"
)

func TestNewConcurrentLimit(t *testing.T) {
	assert := assert.New(t)
	fn := NewConcurrentLimit([]string{
		"q:key",
	}, 10*time.Millisecond, "TestNewConcurrentLimit")
	key := util.RandomString(8)
	req := httptest.NewRequest("GET", "/?key="+key, nil)
	c := elton.NewContext(nil, req)
	c.Next = func() error {
		return nil
	}

	// 第一次时，未有相同的key，正常执行
	err := fn(c)
	assert.Nil(err)
	// 有相同的key，无法正常执行
	err = fn(c)
	assert.Equal("statusCode=400, category=elton-concurrent-limiter, message=submit too frequently", err.Error())

	// 锁过期后，可以正常执行
	time.Sleep(20 * time.Millisecond)
	err = fn(c)
	assert.Nil(err)
}

func TestNewConcurrentLimitWithDone(t *testing.T) {
	assert := assert.New(t)
	fn := NewConcurrentLimitWithDone([]string{
		"q:key",
	}, 20*time.Millisecond, "TestNewConcurrentLimitWithDone")
	key := util.RandomString(8)
	req := httptest.NewRequest("GET", "/?key="+key, nil)
	c := elton.NewContext(nil, req)
	// 由于后续有另外的goroutine读取，因此直接先获取一次query
	_ = c.Query()
	c.Next = func() error {
		// 延时响应，方便测试并发访问
		time.Sleep(10 * time.Millisecond)
		return nil
	}
	go func() {
		time.Sleep(2 * time.Millisecond)
		// 由于上一次的并发访问未完成，因此会出错
		err := fn(c)
		assert.Equal("statusCode=400, category=elton-concurrent-limiter, message=submit too frequently", err.Error())
	}()
	err := fn(c)
	assert.Nil(err)

	// 上一次的已完成，可以继续正常执行
	err = fn(c)
	assert.Nil(err)
}

func TestNewIPLimit(t *testing.T) {
	assert := assert.New(t)
	fn := NewIPLimit(1, 5*time.Millisecond, "TestNewIPLimit")
	req := httptest.NewRequest("GET", "/", nil)
	c := elton.NewContext(nil, req)
	c.Next = func() error {
		return nil
	}

	err := fn(c)
	assert.Nil(err)

	// 第二次访问时，则拦截
	err = fn(c)
	assert.Equal("请求过于频繁，请稍候再试！(2/1)", err.(*hes.Error).Message)

	// 等待过期后可正常执行
	time.Sleep(10 * time.Millisecond)
	err = fn(c)
	assert.Nil(err)
}

func TestNewErrorLimit(t *testing.T) {
	assert := assert.New(t)
	fn := NewErrorLimit(1, 5*time.Millisecond, func(c *elton.Context) string {
		return ""
	})
	c := elton.NewContext(nil, httptest.NewRequest("GET", "/", nil))
	customErr := errors.New("abc")
	c.Next = func() error {
		return customErr
	}
	err := fn(c)
	assert.Equal(err, customErr)

	// 第二次执行时，被拦截
	err = fn(c)
	assert.Equal("请求过于频繁，请稍候再试！(1/1)", err.(*hes.Error).Message)
}
