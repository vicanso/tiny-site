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

package middleware

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
	"github.com/vicanso/tiny-site/service"
	"github.com/vicanso/hes"
)

func TestWaitFor(t *testing.T) {
	assert := assert.New(t)
	t.Run("only err occurred", func(t *testing.T) {
		fn := WaitFor(time.Second, true)
		c := elton.NewContext(nil, nil)

		// 有错误发生
		customErr := errors.New("custom error")
		c.Next = func() error {
			return customErr
		}
		start := time.Now()
		err := fn(c)
		assert.Equal(customErr, err)
		assert.LessOrEqual(time.Second.Milliseconds(), time.Since(start).Milliseconds())

		// 无错误发生
		c.Next = func() error {
			return nil
		}
		start = time.Now()
		err = fn(c)
		assert.Nil(err)
		assert.Greater(time.Second.Milliseconds(), time.Since(start).Milliseconds())
	})

	t.Run("wait for all", func(t *testing.T) {
		fn := WaitFor(time.Second)
		c := elton.NewContext(nil, nil)

		// 有错误发生
		customErr := errors.New("custom error")
		c.Next = func() error {
			return customErr
		}
		start := time.Now()
		err := fn(c)
		assert.Equal(customErr, err)
		assert.LessOrEqual(time.Second.Milliseconds(), time.Since(start).Milliseconds())

		// 无错误发生
		c.Next = func() error {
			return nil
		}
		start = time.Now()
		err = fn(c)
		assert.Nil(err)
		assert.LessOrEqual(time.Second.Milliseconds(), time.Since(start).Milliseconds())
	})
}

func TestValidateCaptcha(t *testing.T) {
	assert := assert.New(t)

	magicalCaptcha := "12345"
	fn := ValidateCaptcha(magicalCaptcha)

	req := httptest.NewRequest("GET", "/", nil)
	c := elton.NewContext(nil, req)

	// 请求头未设置
	err := fn(c)
	assert.Equal("图形验证码参数不能为空", err.(*hes.Error).Message)

	// 错误的请求头设置
	req.Header.Set(xCaptchaHeader, "abc")
	err = fn(c)
	assert.Equal("图形验证码参数长度异常(1)", err.(*hes.Error).Message)

	// macical captcha
	req.Header.Set(xCaptchaHeader, "ax:"+magicalCaptcha)
	done := false
	c.Next = func() error {
		done = true
		return nil
	}
	err = fn(c)
	assert.Nil(err)
	assert.True(done)

	// 正确的校验
	info, err := service.GetCaptcha(context.TODO(), "0,0,0", "0,0,0")
	assert.Nil(err)
	done = false
	c.Next = func() error {
		done = true
		return nil
	}
	req.Header.Set(xCaptchaHeader, info.ID+":"+info.Value)
	err = fn(c)
	assert.Nil(err)
	assert.True(done)
}

func TestNewNoCacheWithCondition(t *testing.T) {
	assert := assert.New(t)
	fn := NewNoCacheWithCondition("cache-control", "no-cache")

	// 设置为no-cache
	req := httptest.NewRequest("GET", "/?cache-control=no-cache", nil)
	c := elton.NewContext(httptest.NewRecorder(), req)
	c.Next = func() error {
		c.CacheMaxAge(time.Minute)
		return nil
	}
	err := fn(c)
	assert.Nil(err)
	assert.Equal("no-cache", c.GetHeader("Cache-Control"))

	// 参数不符合，不设置为no-cache
	req = httptest.NewRequest("GET", "/?cache-control=xxx", nil)
	c = elton.NewContext(httptest.NewRecorder(), req)
	c.Next = func() error {
		c.CacheMaxAge(time.Minute)
		return nil
	}
	err = fn(c)
	assert.Nil(err)
	assert.Equal("public, max-age=60", c.GetHeader("Cache-Control"))
}

func TestNewNotFoundHandler(t *testing.T) {
	assert := assert.New(t)

	fn := NewNotFoundHandler()
	req := httptest.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	fn(resp, req)
	assert.Equal(404, resp.Code)
	assert.Equal(elton.MIMEApplicationJSON, resp.Header().Get(elton.HeaderContentType))
	assert.Equal(`{"statusCode":404,"category":"defaultNotFound","message":"Not Found"}`, resp.Body.String())
}

func TestNewMethodNotAllowedHandler(t *testing.T) {
	assert := assert.New(t)

	fn := NewMethodNotAllowedHandler()
	req := httptest.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	fn(resp, req)
	assert.Equal(405, resp.Code)
	assert.Equal(elton.MIMEApplicationJSON, resp.Header().Get(elton.HeaderContentType))
	assert.Equal(`{"statusCode":405,"category":"defaultMethodNotAllowed","message":"Method Not Allowed"}`, resp.Body.String())
}
