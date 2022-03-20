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
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/vicanso/elton"
)

// NewPrefixHandler create a new request url prefix handler
func NewPrefixHandler(prefixes []string) elton.PreHandler {
	var prefixReg *regexp.Regexp
	if len(prefixes) != 0 {
		prefixReg = regexp.MustCompile(fmt.Sprintf(`^(%s)`, strings.Join(prefixes, "|")))
	}
	return func(req *http.Request) {
		if prefixReg == nil {
			return
		}
		req.URL.Path = prefixReg.ReplaceAllString(req.URL.Path, "")
	}
}
