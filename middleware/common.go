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
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/vicanso/elton"
	"github.com/vicanso/tiny-site/cs"
	"github.com/vicanso/tiny-site/helper"
	"github.com/vicanso/tiny-site/log"
	"github.com/vicanso/tiny-site/service"
	"github.com/vicanso/hes"
)

const (
	xCaptchaHeader    = "X-Captcha"
	errCommonCategory = "common-validate"
)

// WaitFor 延时响应中间件，设置最少等待多久再响应
func WaitFor(d time.Duration, onlyErrOccurreds ...bool) elton.Handler {
	ns := d.Nanoseconds()
	onlyErrOccurred := false
	if len(onlyErrOccurreds) != 0 {
		onlyErrOccurred = onlyErrOccurreds[0]
	}
	return func(c *elton.Context) error {
		start := time.Now()
		err := c.Next()
		// 如果未出错，而且配置为仅在出错时才等待
		if err == nil && onlyErrOccurred {
			return err
		}
		use := time.Now().UnixNano() - start.UnixNano()
		// 无论成功还是失败都wait for
		if use < ns {
			time.Sleep(time.Duration(ns-use) * time.Nanosecond)
		}
		return err
	}
}

// ValidateCaptcha 图形难码校验
func ValidateCaptcha(magicalCaptcha string) elton.Handler {
	return func(c *elton.Context) error {
		value := c.GetRequestHeader(xCaptchaHeader)
		if value == "" {
			return hes.New("图形验证码参数不能为空", errCommonCategory)
		}
		arr := strings.Split(value, ":")
		if len(arr) != 2 {
			return hes.New(fmt.Sprintf("图形验证码参数长度异常(%d)", len(arr)), errCommonCategory)
		}
		// 如果有配置万能验证码，则判断是否相等
		if magicalCaptcha != "" && arr[1] == magicalCaptcha {
			return c.Next()
		}
		valid, err := service.ValidateCaptcha(c.Context(), arr[0], arr[1])
		if err != nil {
			if helper.RedisIsNilError(err) {
				err = hes.New("图形验证码已过期，请刷新", errCommonCategory)
			}
			return err
		}
		if !valid {
			return hes.New("图形验证码错误", errCommonCategory)
		}
		return c.Next()
	}
}

// NewNoCacheWithCondition 创建no cache的中间件，此中间件根据设置的key value来判断是否设置为no cache
func NewNoCacheWithCondition(key, value string) elton.Handler {
	return func(c *elton.Context) error {
		err := c.Next()
		if c.QueryParam(key) == value {
			c.NoCache()
		}
		return err
	}
}

// NewNotFoundHandler 创建404 not found的处理函数
func NewNotFoundHandler() http.HandlerFunc {
	// 对于404的请求，不会执行中间件，一般都是因为攻击之类才会导致大量出现404，
	notFoundErrBytes := (&hes.Error{
		Message:    "Not Found",
		StatusCode: http.StatusNotFound,
		Category:   "defaultNotFound",
	}).ToJSON()
	return func(resp http.ResponseWriter, req *http.Request) {
		ip := elton.GetClientIP(req)
		log.Info(req.Context()).
			Str("category", "404").
			Str("ip", ip).
			Str("method", req.Method).
			Str("uri", req.RequestURI).
			Msg("")

		status := http.StatusNotFound
		resp.Header().Set(elton.HeaderContentType, elton.MIMEApplicationJSON)
		resp.WriteHeader(status)
		_, err := resp.Write(notFoundErrBytes)
		if err != nil {
			log.Error(req.Context()).
				Str("ip", ip).
				Str("uri", req.RequestURI).
				Err(err).
				Msg("404 response fail")
		}

		tags := map[string]string{
			cs.TagMethod: req.Method,
			cs.TagRoute:  "404",
		}
		fields := map[string]interface{}{
			cs.FieldIP:     ip,
			cs.FieldURI:    req.RequestURI,
			cs.FieldStatus: status,
		}
		helper.GetInfluxDB().Write(cs.MeasurementHTTPStats, tags, fields)
	}
}

// NewMethodNotAllowedHandler 创建method not allowed的处理函数
func NewMethodNotAllowedHandler() http.HandlerFunc {
	methodNotAllowedErrBytes := (&hes.Error{
		Message:    "Method Not Allowed",
		StatusCode: http.StatusMethodNotAllowed,
		Category:   "defaultMethodNotAllowed",
	}).ToJSON()
	return func(resp http.ResponseWriter, req *http.Request) {
		ip := elton.GetClientIP(req)
		log.Info(req.Context()).
			Str("category", "405").
			Str("ip", ip).
			Str("method", req.Method).
			Str("uri", req.RequestURI).
			Msg("")
		resp.Header().Set(elton.HeaderContentType, elton.MIMEApplicationJSON)
		status := http.StatusMethodNotAllowed
		resp.WriteHeader(status)
		_, err := resp.Write(methodNotAllowedErrBytes)
		if err != nil {
			log.Error(req.Context()).
				Str("ip", ip).
				Str("uri", req.RequestURI).
				Err(err).
				Msg("405 response fail")
		}
		tags := map[string]string{
			cs.TagMethod: req.Method,
			cs.TagRoute:  "405",
		}
		fields := map[string]interface{}{
			cs.FieldIP:     ip,
			cs.FieldURI:    req.RequestURI,
			cs.FieldStatus: status,
		}
		helper.GetInfluxDB().Write(cs.MeasurementHTTPStats, tags, fields)
	}
}
