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
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
	routermock "github.com/vicanso/tiny-site/router_mock"
)

func TestNewRouterMocker(t *testing.T) {
	assert := assert.New(t)
	routeConfigs := []*routermock.RouterMock{
		{
			Method:   "GET",
			Route:    "/route-match",
			Response: "route match response",
			Status:   200,
		},
		{
			Method:     "GET",
			Route:      "/route-uri-match",
			URL:        "/route-uri-match?a=1",
			Response:   "route uri match response",
			CotentType: "text/html",
			Status:     201,
		},
	}
	getConfig := func(method, route string) (found *routermock.RouterMock) {
		for _, r := range routeConfigs {
			if r.Method == method && r.Route == route {
				found = r
				break
			}
		}
		return
	}
	mid := NewRouterMocker(getConfig)

	tests := []struct {
		desc        string
		r           *http.Request
		nextDone    bool
		response    *bytes.Buffer
		status      int
		contentType string
	}{
		{
			desc:     "no match",
			r:        httptest.NewRequest("GET", "/no-match", nil),
			nextDone: true,
		},
		{
			desc:        "route match",
			r:           httptest.NewRequest("GET", "/route-match", nil),
			response:    bytes.NewBufferString("route match response"),
			status:      200,
			contentType: elton.MIMEApplicationJSON,
		},
		{
			desc:     "route match but uri not match",
			r:        httptest.NewRequest("GET", "/route-uri-match?b=1", nil),
			nextDone: true,
		},
		{
			desc:        "route and uril match",
			r:           httptest.NewRequest("GET", "/route-uri-match?a=1", nil),
			response:    bytes.NewBufferString("route uri match response"),
			status:      201,
			contentType: "text/html",
		},
	}
	for _, tt := range tests {

		c := elton.NewContext(httptest.NewRecorder(), tt.r)
		c.Route = tt.r.URL.Path
		done := false
		c.Next = func() error {
			done = true
			return nil
		}
		err := mid(c)
		assert.Nil(err)
		assert.Equal(tt.nextDone, done)
		assert.Equal(tt.response, c.BodyBuffer)
		assert.Equal(tt.status, c.StatusCode)
		assert.Equal(tt.contentType, c.GetHeader(elton.HeaderContentType))
	}
}
