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

package interceptor

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/dop251/goja"
	"github.com/vicanso/elton"
	"github.com/vicanso/tiny-site/asset"
	"github.com/vicanso/tiny-site/util"
	"go.uber.org/atomic"
)

type httpInterceptorScript struct {
	Before string
	After  string
}
type httpInterceptors struct {
	scripts    *atomic.Value
	baseScript string
}

func (interceptors *httpInterceptors) getScripts() map[string]*httpInterceptorScript {
	value := interceptors.scripts.Load()
	if value == nil {
		return nil
	}
	scripts, ok := value.(map[string]*httpInterceptorScript)
	if !ok {
		return nil
	}
	return scripts
}

func (interceptors *httpInterceptors) Get(router string) *httpInterceptorScript {
	scripts := interceptors.getScripts()
	if scripts == nil {
		return nil
	}
	return scripts[router]
}

func UpdateHTTPInterceptors(arr []string) {
	scripts := make(map[string]*httpInterceptorScript)
	for _, item := range arr {
		script := make(map[string]string)
		_ = json.Unmarshal([]byte(item), &script)
		router := script["router"]
		// 只根据是否有router来判断是否正确
		if router == "" {
			continue
		}
		scripts[router] = &httpInterceptorScript{
			Before: script["before"],
			After:  script["after"],
		}
	}
	currentHTTPInterceptors.scripts.Store(scripts)
}

// http服务器接收请求
// 其中query为了简便直接使用了map[string]string替换map[string][]string
type httpServerRequest struct {
	URI   string            `json:"uri"`
	Query map[string]string `json:"query"`
	// 是否有修改query
	ModifiedQuery bool                   `json:"modifiedQuery"`
	Body          map[string]interface{} `json:"body"`
	// 是否有修改body
	ModifiedBody bool `json:"modifiedBody"`
}

// http服务响应数据
type httpServerResponse struct {
	// 响应状态码
	Status int `json:"status"`
	// 响应头
	Header map[string]string `json:"header"`
	// 响应数据
	Body map[string]interface{} `json:"body"`
}

type httpServerInterceptor struct {
	Before func() (*httpServerResponse, error)
	After  func() (*httpServerResponse, error)
}

var currentHTTPInterceptors = newHTTPInterceptors()

func newHTTPInterceptors() *httpInterceptors {
	script, _ := asset.GetFS().ReadFile("http_server_interceptor.js")
	return &httpInterceptors{
		scripts:    &atomic.Value{},
		baseScript: string(script),
	}
}

func newScript(script string) string {
	return fmt.Sprintf(`%s;(function() {
		%s
	})();
	`, currentHTTPInterceptors.baseScript, script)
}

func newHTTPServerRequest(c *elton.Context) *httpServerRequest {
	query := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		query[key] = values[0]
	}
	req := &httpServerRequest{
		URI:   c.Request.URL.RequestURI(),
		Query: query,
		Body:  make(map[string]interface{}),
	}
	if len(c.RequestBody) != 0 {
		_ = json.Unmarshal(c.RequestBody, &req.Body)
	}
	return req
}

func newHTTPServerResponse() *httpServerResponse {
	return &httpServerResponse{
		Header: make(map[string]string),
		Body:   make(map[string]interface{}),
	}
}

// 设置响应数据
func (resp *httpServerResponse) SetResponse(c *elton.Context) {
	c.StatusCode = resp.Status
	for k, v := range resp.Header {
		c.SetHeader(k, v)
	}
	buf, _ := json.Marshal(resp.Body)
	c.BodyBuffer = bytes.NewBuffer(buf)
}

func NewHTTPServer(c *elton.Context) (inter *httpServerInterceptor, err error) {
	script := currentHTTPInterceptors.Get(c.Request.Method + " " + c.Route)
	if script == nil {
		return nil, nil
	}

	vm := goja.New()

	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	req := newHTTPServerRequest(c)
	err = vm.Set("req", req)
	if err != nil {
		return
	}
	resp := newHTTPServerResponse()
	err = vm.Set("resp", resp)
	if err != nil {
		return
	}

	inter = &httpServerInterceptor{
		Before: func() (*httpServerResponse, error) {
			if script.Before == "" {
				return nil, nil
			}
			_, err := vm.RunString(newScript(script.Before))
			if err != nil {
				return nil, err
			}
			// 如果有修改query，则重新生成
			if req.ModifiedQuery {
				query := c.Request.URL.Query()
				for k, v := range req.Query {
					query.Set(k, v)
				}
				c.Request.URL.RawQuery = query.Encode()
				c.Request.RequestURI = c.Request.URL.RequestURI()
			}
			// 如果有修改body
			if req.ModifiedBody {
				body := util.MergeMapStringInterface(req.Body, req.Body)
				c.RequestBody, _ = json.Marshal(body)
			}
			return resp, nil
		},
		After: func() (*httpServerResponse, error) {
			if script.After == "" {
				return nil, nil
			}
			_ = json.Unmarshal(c.BodyBuffer.Bytes(), &resp.Body)
			_, err := vm.RunString(newScript(script.After))
			if err != nil {
				return nil, err
			}
			return resp, nil
		},
	}
	return
}
