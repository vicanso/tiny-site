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

package controller

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
	"github.com/vicanso/tiny-site/service"
	"github.com/vicanso/go-axios"
)

func TestCommonCtrl(t *testing.T) {
	ctrl := commonCtrl{}
	t.Run("ping", func(t *testing.T) {
		assert := assert.New(t)
		c := elton.NewContext(nil, nil)
		err := ctrl.ping(c)
		assert.Nil(err)
		assert.Equal("pong", c.BodyBuffer.String())
	})

	t.Run("location", func(t *testing.T) {
		ins := service.LocationIns
		originalData := []byte(`{"ip":"1.1.1.1","country":"澳大利亚"}`)
		done := ins.Mock(&axios.Response{
			Data:   originalData,
			Status: 200,
		})
		defer done()

		assert := assert.New(t)
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set(elton.HeaderXForwardedFor, "1.1.1.1")
		c := elton.NewContext(nil, req)
		err := ctrl.location(c)
		assert.Nil(err)
		assert.NotNil(c.Body)
		buf, _ := json.Marshal(c.Body)
		assert.Equal(originalData, buf)
	})

	t.Run("randomKeys", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest("GET", "/?n=10&size=1", nil)
		c := elton.NewContext(nil, req)

		err := ctrl.randomKeys(c)
		assert.Nil(err)
		assert.NotNil(c.Body)
		m, ok := c.Body.(map[string][]string)
		assert.True(ok)
		assert.Equal(1, len(m["keys"]))
		assert.Equal(10, len(m["keys"][0]))
	})

	t.Run("captcha", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest("GET", "/", nil)
		c := elton.NewContext(nil, req)
		err := ctrl.captcha(c)
		assert.Nil(err)
		_, ok := c.Body.(*service.CaptchaInfo)
		assert.True(ok)
	})
}
