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

func TestNewEntry(t *testing.T) {
	assert := assert.New(t)
	req := httptest.NewRequest("GET", "/", nil)
	c := elton.NewContext(httptest.NewRecorder(), req)
	c.ID = "abc"
	done := false
	c.Next = func() error {
		done = true
		return nil
	}
	doneEntry := false
	doneExit := false

	fn := NewEntry(func() int32 {
		doneEntry = true
		return 0
	}, func() int32 {
		doneExit = true
		return 0
	})
	err := fn(c)
	assert.Nil(err)
	assert.True(done)
	assert.True(doneEntry)
	assert.True(doneExit)
	assert.Equal("abc", c.GetHeader(xResponseID))
	assert.Equal("no-cache", c.GetHeader("Cache-Control"))
}
