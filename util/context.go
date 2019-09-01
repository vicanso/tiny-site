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

package util

import (
	"github.com/vicanso/elton"

	"github.com/vicanso/tiny-site/config"
)

// GetTrackID get track id
func GetTrackID(c *elton.Context) string {
	trackCookie := config.GetTrackKey()
	if trackCookie == "" {
		return ""
	}
	cookie, _ := c.Cookie(trackCookie)
	if cookie == nil {
		return ""
	}
	return cookie.Value
}

// GetSessionID get session id
func GetSessionID(c *elton.Context) string {
	scf := config.GetSessionConfig()
	cookie, _ := c.Cookie(scf.Key)
	if cookie == nil {
		return ""
	}
	return cookie.Value
}
