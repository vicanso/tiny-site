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

package request

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/tidwall/gjson"
	"github.com/vicanso/tiny-site/cs"
	"github.com/vicanso/tiny-site/helper"
	"github.com/vicanso/tiny-site/log"
	"github.com/vicanso/tiny-site/util"
	"github.com/vicanso/go-axios"
	"github.com/vicanso/hes"
)

func newOnDone(serviceName string) axios.OnDone {
	return func(conf *axios.Config, resp *axios.Response, err error) {
		stats := axios.GetStats(conf, err)
		ht := conf.HTTPTrace

		use := ""

		tags := map[string]string{
			cs.TagService: serviceName,
			cs.TagRoute:   stats.Route,
			cs.TagMethod:  stats.Method,
			cs.TagResult:  strconv.Itoa(stats.Result),
		}
		fields := map[string]interface{}{
			cs.FieldURI:    stats.URI,
			cs.FieldStatus: stats.Status,
		}
		if ht != nil {
			use = ht.Stats().String()
			fields[cs.FieldReused] = stats.Reused
			fields[cs.FieldLatency] = stats.Use
			if stats.DNSUse != 0 {
				fields[cs.FieldDNSUse] = stats.DNSUse
			}
			if stats.TCPUse != 0 {
				fields[cs.FieldTCPUse] = stats.TCPUse
			}
			if stats.TLSUse != 0 {
				fields[cs.FieldTLSUse] = stats.TLSUse
			}
			if stats.ServerProcessingUse != 0 {
				fields[cs.FieldProcessingUse] = stats.ServerProcessingUse
			}
			if stats.ContentTransferUse != 0 {
				fields[cs.FieldTransferUse] = stats.ContentTransferUse
			}
			fields[cs.FieldAddr] = stats.Addr
		}
		message := ""
		if err != nil {
			he := hes.Wrap(err)
			message = he.Error()
			fields[cs.FieldError] = message
			errCategory := he.Category
			if errCategory != "" {
				fields[cs.FieldErrCategory] = errCategory
			}
			if he.Exception {
				fields[cs.FieldException] = true
			}
		}
		// 输出响应数据，如果响应数据为隐私数据可不输出
		var data interface{}
		size := stats.Size
		if resp != nil {
			data = resp.UnmarshalData
		}
		// 由于http请求是较频繁的操作，因此判断是否启用debug再输出
		if log.DebugEnabled() {
			respData := ""
			if resp != nil {
				respData = string(resp.Data)
			}
			log.Debug(conf.Context).
				Str("curl", conf.CURL()).
				Str("data", respData).
				Msg("request log")
		}
		requestURL := stats.URI
		urlInfo, _ := url.Parse(requestURL)
		if urlInfo != nil {
			requestURL = urlInfo.RequestURI()
		}
		event := log.Info(conf.Context).
			Str("category", "requestStats").
			Str("service", serviceName).
			Str("method", stats.Method).
			Str("route", stats.Route).
			Str("url", requestURL)
		if len(conf.Params) != 0 {
			event = event.Dict("params", log.Struct(conf.Params))
		}
		if len(conf.Query) != 0 {
			event = event.Dict("query", log.URLValues(conf.Query))
		}
		if data != nil {
			event = event.Dict("data", log.Struct(data))
		} else if resp != nil {
			event = event.Str("data", string(resp.Data))
		}
		event.Int("size", size).
			Int("status", stats.Status).
			Str("addr", stats.Addr).
			Bool("reused", stats.Reused).
			Str("use", use)
		if message != "" {
			event = event.Str("error", message)
		}
		event.Msg("")
		helper.GetInfluxDB().Write(cs.MeasurementHTTPRequest, tags, fields)
	}
}

// newConvertResponseToError 将http响应码为>=400的转换为出错
func newConvertResponseToError() axios.ResponseInterceptor {
	return func(resp *axios.Response) error {
		if resp.Status >= 400 {
			message := gjson.GetBytes(resp.Data, "message").String()
			exception := false
			if message == "" {
				message = util.CutRune(string(resp.Data), 30)
				// 如果出错响应不符合，则认为是异常响应
				exception = true
			}
			he := hes.NewWithStatusCode(message, resp.Status)
			he.Exception = exception
			return he
		}
		return nil
	}
}

// newOnError 新建error的处理函数
func newOnError(serviceName string) axios.OnError {
	return func(err error, conf *axios.Config) error {
		code := -1
		if conf.Response != nil {
			code = conf.Response.Status
		}
		he := hes.Wrap(err)
		if code >= http.StatusBadRequest {
			he.StatusCode = code
		}
		// 如果状态码>500，认为是异常（一般nginx之类才会使用502 503之类的状态码）
		if code > http.StatusInternalServerError {
			he.Exception = true
		}
		// 如果未设置http响应码(<400)，则设置为500
		if he.StatusCode < http.StatusBadRequest {
			he.StatusCode = http.StatusInternalServerError
		}

		// 如果为空，则通过error获取
		if he.Category == "" {
			he.Category = axios.GetInternalErrorCategory(err)
			// 如果返回错误类型，则认为异常
			// 因为返回出错均是网络连接上的异常
			if he.Category != "" {
				he.Exception = true
			}
		}

		if !util.IsProduction() {
			he.AddExtra("requestRoute", conf.Route)
			he.AddExtra("requestService", serviceName)
			he.AddExtra("requestCURL", conf.CURL())
		}
		return he
	}
}
