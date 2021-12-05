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

package controller

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
	"github.com/vicanso/tiny-site/service"
	"github.com/vicanso/hes"
)

func TestCommonCtrl(t *testing.T) {
	assert := assert.New(t)
	ctrl := commonCtrl{}

	t.Run("ping", func(t *testing.T) {
		originalStatus := service.GetApplicationStatus()
		defer service.SetApplicationStatus(originalStatus)

		service.SetApplicationStatus(service.ApplicationStatusStopped)
		c := elton.NewContext(httptest.NewRecorder(), nil)
		err := ctrl.ping(c)
		assert.Equal("应用服务不可用", err.(*hes.Error).Message)

		service.SetApplicationStatus(service.ApplicationStatusRunning)
		err = ctrl.ping(c)
		assert.Nil(err)
		assert.Equal("pong", c.BodyBuffer.String())
	})

	t.Run("getApplicationInfo", func(t *testing.T) {
		c := elton.NewContext(httptest.NewRecorder(), nil)
		err := ctrl.getApplicationInfo(c)
		assert.Nil(err)
		info, ok := c.Body.(*applicationInfoResp)
		assert.True(ok)
		assert.True(strings.HasPrefix(info.GO, "go1.1"))
		assert.Equal("public, max-age=60", c.Header().Get(elton.HeaderCacheControl))
	})

	t.Run("getCaptcha", func(t *testing.T) {
		c := elton.NewContext(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
		err := ctrl.getCaptcha(c)
		assert.Nil(err)
		assert.Equal("no-store", c.Header().Get(elton.HeaderCacheControl))
		_, ok := c.Body.(*service.CaptchaInfo)
		assert.True(ok)
	})

	t.Run("getPerformance", func(t *testing.T) {
		c := elton.NewContext(nil, httptest.NewRequest("GET", "/", nil))
		err := ctrl.getPerformance(c)
		assert.Nil(err)
		_, ok := c.Body.(*service.Performance)
		assert.True(ok)
	})

	t.Run("listStatus", func(t *testing.T) {
		c := elton.NewContext(httptest.NewRecorder(), nil)
		err := ctrl.listStatus(c)
		assert.Nil(err)
		data, ok := c.Body.(*statusListResp)
		assert.True(ok)
		assert.NotEmpty(data.Statuses)
		assert.Equal("public, max-age=300", c.GetHeader(elton.HeaderCacheControl))
	})

	t.Run("getRandomKeys", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		c := elton.NewContext(nil, req)
		err := ctrl.getRandomKeys(c)
		assert.Nil(err)
		data, ok := c.Body.(*randomKeysResp)
		assert.True(ok)
		assert.Equal(1, len(data.Keys))
	})
}
