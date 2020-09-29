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
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/vicanso/elton"

	"github.com/vicanso/hes"

	"github.com/vicanso/go-axios"
	"github.com/vicanso/tiny-site/cs"
	"github.com/vicanso/tiny-site/util"
	"go.uber.org/zap"
)

// newHTTPStats http stats
func newHTTPStats(serviceName string) axios.ResponseInterceptor {
	return func(resp *axios.Response) (err error) {
		conf := resp.Config

		ht := conf.HTTPTrace

		reused := false
		addr := ""
		use := ""
		if ht != nil {
			reused = ht.Reused
			addr = ht.Addr
			use = ht.Stats().Total.String()
		}
		// TODO 统计可以写入influxdb
		logger.Info("http stats",
			zap.String("service", serviceName),
			zap.String("cid", conf.GetString(cs.CID)),
			zap.String("method", conf.Method),
			zap.String("route", conf.Route),
			zap.String("url", conf.URL),
			zap.Int("status", resp.Status),
			zap.String("addr", addr),
			zap.Bool("reused", reused),
			zap.String("use", use),
		)
		return
	}
}

// newConvertResponseToError convert http response(4xx, 5xx) to error
func newConvertResponseToError(serviceName string) axios.ResponseInterceptor {
	return func(resp *axios.Response) (err error) {
		if resp.Status >= 400 {
			message := standardJSON.Get(resp.Data, "message").ToString()
			if message == "" {
				message = string(resp.Data)
			}
			err = errors.New(message)
		}
		return
	}
}

// newOnError 新建error的处理函数
func newOnError(serviceName string) axios.OnError {
	return func(err error, conf *axios.Config) (newErr error) {
		id := conf.GetString(cs.CID)
		code := -1
		if conf.Response != nil {
			code = conf.Response.Status
		}
		var he *hes.Error
		if hes.IsError(err) {
			he, _ = err.(*hes.Error)
		} else {
			he = hes.Wrap(err)
			he.StatusCode = code
		}
		// 如果未设置http响应码，则设置为500
		if he.StatusCode == 0 {
			he.StatusCode = http.StatusInternalServerError
		}

		if he.Extra == nil {
			he.Extra = make(map[string]interface{})
		}

		// 请求超时
		e, ok := err.(*url.Error)
		if ok && e.Timeout() {
			he.Extra["category"] = "timeout"
		}
		if !util.IsProduction() {
			he.Extra["requestRoute"] = conf.Route
			he.Extra["requestService"] = serviceName
			he.Extra["requestCURL"] = conf.CURL()
			// TODO 是否非生产环境增加更多的信息，方便测试时确认问题
		}
		newErr = he
		logger.Info("http error",
			zap.String("service", serviceName),
			zap.String("cid", id),
			zap.String("method", conf.Method),
			zap.String("url", conf.URL),
			zap.String("error", err.Error()),
		)
		return
	}
}

// NewInstance new an instance
func NewInstance(serviceName, baseURL string, timeout time.Duration) *axios.Instance {
	return axios.NewInstance(&axios.InstanceConfig{
		EnableTrace: true,
		Timeout:     timeout,
		OnError:     newOnError(serviceName),
		BaseURL:     baseURL,
		ResponseInterceptors: []axios.ResponseInterceptor{
			newHTTPStats(serviceName),
			newConvertResponseToError(serviceName),
		},
	})
}

// AttachWithContext attach with context
func AttachWithContext(conf *axios.Config, c *elton.Context) {
	if c == nil || conf == nil {
		return
	}
	conf.Set(cs.CID, c.ID)
}
