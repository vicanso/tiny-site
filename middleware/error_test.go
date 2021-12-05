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
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
	"github.com/vicanso/hes"
)

func TestNewError(t *testing.T) {
	assert := assert.New(t)

	fn := NewError()
	req := httptest.NewRequest("GET", "/", nil)
	c := elton.NewContext(nil, req)
	c.Route = "/"
	c.Next = func() error {
		return nil
	}
	err := fn(c)
	assert.Nil(err)
	assert.Empty(c.BodyBuffer)

	// 非hes的认为unexpected error
	c.Next = func() error {
		return errors.New("abc")
	}
	err = fn(c)
	assert.Nil(err)
	assert.Equal(500, c.StatusCode)
	assert.Equal(`{"statusCode":500,"message":"abc","exception":true,"extra":{"route":"/"}}`, c.BodyBuffer.String())

	// 自定义的hes出错
	c.Next = func() error {
		return hes.New("abc")
	}
	err = fn(c)
	assert.Nil(err)
	assert.Equal(400, c.StatusCode)
	assert.Equal(`{"statusCode":400,"message":"abc","extra":{"route":"/"}}`, c.BodyBuffer.String())
}
