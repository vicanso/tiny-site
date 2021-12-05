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
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
)

func TestNewIPBlocker(t *testing.T) {
	assert := assert.New(t)
	blockFn := func(ip string) bool {
		return ip == "1.1.1.1"
	}
	fn := NewIPBlocker(blockFn)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set(elton.HeaderXForwardedFor, "1.1.1.1")
	c := elton.NewContext(nil, req)
	err := fn(c)
	assert.Equal(ErrIPNotAllow, err)

	req.Header.Del(elton.HeaderXForwardedFor)
	// 由于context的ip会缓存，因此重新创建
	c = elton.NewContext(nil, req)
	done := false
	c.Next = func() error {
		done = true
		return nil
	}
	err = fn(c)
	assert.Nil(err)
	assert.True(done)
}
