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
	"github.com/vicanso/elton"
	"github.com/vicanso/tiny-site/service"
	"github.com/vicanso/tiny-site/util"
)

const (
	xResponseID = "X-Response-Id"
)

type EntryFunc func() int32
type ExitFunc func() int32

// NewEntry create an entry middleware
func NewEntry(entryFn EntryFunc, exitFn ExitFunc) elton.Handler {
	return func(c *elton.Context) error {
		// 如果请求头指定connection: close
		// server conn stats中无法判断现在处理请求事件
		// 处理完成时，需要手工将现在处理的连接数-1
		if c.GetRequestHeader("Connection") == "close" {
			defer service.DecHTTPConnProcessing()
		}
		entryFn()
		defer exitFn()
		if c.ID != "" {
			c.SetHeader(xResponseID, c.ID)
			c.WithContext(util.SetTraceID(c.Context(), c.ID))
		}
		deviceID := util.GetTrackID(c)
		if deviceID != "" {
			c.WithContext(util.SetDeviceID(c.Context(), deviceID))
		}

		// 测试环境返回x-forwarded-for，方便确认链路
		if !util.IsProduction() {
			c.SetHeader(elton.HeaderXForwardedFor, c.GetRequestHeader(elton.HeaderXForwardedFor))
		}

		// 设置所有的请求响应默认都为no cache
		c.NoCache()

		return c.Next()
	}
}
