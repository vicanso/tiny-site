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
	"net/http"
	"strings"
	"time"

	"github.com/vicanso/elton"
	"github.com/vicanso/tiny-site/service"
	"github.com/vicanso/hes"
)

const (
	xCaptchHeader = "X-Captcha"
	errCategory   = "common-validate"
)

var (
	errQueryNotAllow = &hes.Error{
		StatusCode: http.StatusBadRequest,
		Message:    "query is not allowed",
		Category:   errCategory,
	}
	errCaptchIsInvalid = &hes.Error{
		StatusCode: http.StatusBadRequest,
		Message:    "captcha is invalid",
		Category:   errCategory,
	}
)

// NoQuery no query middleware
func NoQuery(c *elton.Context) (err error) {
	if c.Request.URL.RawQuery != "" {
		err = errQueryNotAllow
		return
	}
	return c.Next()
}

// WaitFor at least wait for duration
func WaitFor(d time.Duration) elton.Handler {
	ns := d.Nanoseconds()
	return func(c *elton.Context) (err error) {
		start := time.Now()
		err = c.Next()
		use := time.Now().UnixNano() - start.UnixNano()
		// 无论成功还是失败都wait for
		if use < ns {
			time.Sleep(time.Duration(ns-use) * time.Nanosecond)
		}
		return
	}
}

// ValidateCaptch validate chapter
func ValidateCaptch() elton.Handler {
	return func(c *elton.Context) (err error) {
		value := c.GetRequestHeader(xCaptchHeader)
		if value == "" {
			err = errCaptchIsInvalid
			return
		}
		arr := strings.Split(value, ":")
		if len(arr) != 2 {
			err = errCaptchIsInvalid
			return
		}
		valid, err := service.ValidateCaptcha(arr[0], arr[1])
		if err != nil {
			return err
		}
		if !valid {
			err = errCaptchIsInvalid
			return
		}
		return c.Next()
	}
}
