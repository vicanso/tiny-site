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
	"time"

	"github.com/vicanso/elton"

	"github.com/vicanso/hes"

	"github.com/vicanso/tiny-site/cs"
	"github.com/vicanso/go-axios"
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

// newOnError new an error listener
func newOnError(serviceName string) axios.OnError {
	return func(err error, conf *axios.Config) (newErr error) {
		e, ok := err.(*axios.Error)
		id := conf.GetString(cs.CID)
		if ok {
			code := e.Code

			he := &hes.Error{
				StatusCode: code,
				Message:    e.Message,
				ID:         id,
			}
			if code < http.StatusBadRequest {
				he.Exception = true
				he.StatusCode = http.StatusInternalServerError
			}

			// 请求超时
			if e.Timeout() {
				he.Message = "Timeout"
			}
			if !isProduction() {
				he.Extra = map[string]interface{}{
					"route":   conf.Route,
					"service": serviceName,
				}
			}
			newErr = he
		}
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
