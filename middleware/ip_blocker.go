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
	"net/http"

	"github.com/vicanso/elton"
	"github.com/vicanso/hes"
)

var (
	ErrIPNotAllow = &hes.Error{
		StatusCode: http.StatusBadRequest,
		Message:    "request is forbidden",
		Category:   "IB",
	}
)

type IPBlockFunc func(string) bool

// NewIPBlocker create a new block ip middleware
func NewIPBlocker(fn IPBlockFunc) elton.Handler {
	return func(c *elton.Context) error {
		if fn(c.RealIP()) {
			return ErrIPNotAllow
		}
		return c.Next()
	}
}
