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
package util

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
	"github.com/vicanso/tiny-site/config"
)

func TestGetTrackID(t *testing.T) {
	assert := assert.New(t)
	sessionConfig = config.MustGetSessionConfig()
	req := httptest.NewRequest("GET", "/", nil)
	cookie := http.Cookie{
		Name:  sessionConfig.TrackKey,
		Value: "abcd",
	}
	req.AddCookie(&cookie)
	c := elton.NewContext(nil, req)
	assert.Equal(cookie.Value, GetTrackID(c))
}

func TestGetSessionID(t *testing.T) {
	assert := assert.New(t)
	sessionConfig = config.MustGetSessionConfig()
	req := httptest.NewRequest("GET", "/", nil)
	cookie := http.Cookie{
		Name:  sessionConfig.Key,
		Value: "abcd",
	}
	req.AddCookie(&cookie)
	c := elton.NewContext(nil, req)
	assert.Equal(cookie.Value, GetSessionID(c))
}
