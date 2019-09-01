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
	"bytes"

	"github.com/vicanso/elton"
	"github.com/vicanso/tiny-site/service"
)

// NewRouterController create a router controller
func NewRouterController() elton.Handler {
	return func(c *elton.Context) (err error) {
		routerConfig := service.RouterGetConfig(c.Request.Method, c.Route)
		if routerConfig == nil {
			return c.Next()
		}

		c.StatusCode = routerConfig.Status
		contentType := routerConfig.CotentType
		if contentType == "" {
			contentType = elton.MIMEApplicationJSON
		}
		c.SetHeader(elton.HeaderContentType, contentType)
		c.BodyBuffer = bytes.NewBufferString(routerConfig.Response)
		return nil
	}
}
