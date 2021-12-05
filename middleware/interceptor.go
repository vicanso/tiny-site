// Copyright 2021 tree xie
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
	"github.com/vicanso/tiny-site/interceptor"
)

func NewInterceptor() elton.Handler {
	return func(c *elton.Context) error {
		inter, err := interceptor.NewHTTPServer(c)
		if err != nil {
			return err
		}
		// 如果返回空表示没有设置interceptor
		if inter == nil {
			return c.Next()
		}
		// 前置处理
		resp, err := inter.Before()
		if err != nil {
			return err
		}
		// 如果状态码不为0，则表示已设置响应数据
		if resp != nil && resp.Status != 0 {
			resp.SetResponse(c)
			return nil
		}
		err = c.Next()
		if err != nil {
			return err
		}
		// 后置处理
		resp, err = inter.After()
		if err != nil {
			return err
		}
		// 如果状态码不为0，则表示重设响应数据
		if resp != nil && resp.Status != 0 {
			resp.SetResponse(c)
			return nil
		}
		return nil
	}
}
