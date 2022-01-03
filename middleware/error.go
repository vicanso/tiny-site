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
	"bytes"
	"net/http"

	"github.com/vicanso/elton"
	"github.com/vicanso/hes"
	"github.com/vicanso/tiny-site/cs"
	"github.com/vicanso/tiny-site/helper"
	"github.com/vicanso/tiny-site/log"
	"github.com/vicanso/tiny-site/session"
	"github.com/vicanso/tiny-site/util"
)

// New Error handler
func NewError() elton.Handler {

	return func(c *elton.Context) error {
		err := c.Next()
		if err == nil {
			return nil
		}
		uri := c.Request.RequestURI
		he, ok := err.(*hes.Error)
		if !ok {
			// 如果不是以http error的形式返回的error则为非主动抛出错误
			he = hes.NewWithError(err)
			he.StatusCode = http.StatusInternalServerError
			he.Exception = true
		} else {
			// 避免修改了原有的error对象
			he = he.Clone()
		}
		if he.StatusCode == 0 {
			he.StatusCode = http.StatusInternalServerError
		}
		account := ""
		tid := util.GetDeviceID(c.Context())
		us := session.NewUserSession(c)
		if us != nil && us.IsLogin() {
			account = us.MustGetInfo().Account
		}

		ip := c.RealIP()
		log.Info(c.Context()).
			Str("catgory", "httpError").
			Bool("exception", he.Exception).
			Str("ip", ip).
			Str("method", c.Request.Method).
			Str("route", c.Route).
			Str("uri", uri).
			Str("error", he.Error()).
			Msg("")

		sid := util.GetSessionID(c)

		he.AddExtra("route", c.Route)
		// 记录用户相关信息
		fields := map[string]interface{}{
			cs.FieldStatus:    he.StatusCode,
			cs.FieldError:     he.Error(),
			cs.FieldURI:       uri,
			cs.FieldException: he.Exception,
			cs.FieldIP:        ip,
			cs.FieldSID:       sid,
			cs.FieldTID:       tid,
		}
		if account != "" {
			fields[cs.FieldAccount] = account
		}
		tags := map[string]string{
			cs.TagMethod: c.Request.Method,
			cs.TagRoute:  c.Route,
		}
		if he.Category != "" {
			tags[cs.TagCategory] = he.Category
		}

		helper.GetInfluxDB().Write(cs.MeasurementHTTPError, tags, fields)
		c.StatusCode = he.StatusCode
		c.SetContentTypeByExt(".json")
		c.BodyBuffer = bytes.NewBuffer(he.ToJSON())
		return nil
	}
}
