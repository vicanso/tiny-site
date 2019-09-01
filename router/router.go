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

package router

import (
	"regexp"

	"github.com/vicanso/elton"
	"github.com/vicanso/hes"
)

var (
	// groupList 路由组列表
	groupList = make([]*elton.Group, 0)
)

// NewGroup new router group
func NewGroup(path string, handlerList ...elton.Handler) *elton.Group {
	// 如果配置文件中有配置路由
	g := elton.NewGroup(path, handlerList...)
	groupList = append(groupList, g)
	return g
}

// Init init router
func Init(d *elton.Elton) {
	for _, g := range groupList {
		d.AddGroup(g)
	}

	configIDReg := regexp.MustCompile(`^[1-9][0-9]{0,3}$`)
	d.AddValidator("configID", func(value string) error {
		if !configIDReg.MatchString(value) {
			return hes.New("config id should be numbers")
		}
		return nil
	})

	// 如果用户量增大，需要调整此限制
	userIDReg := regexp.MustCompile(`^[1-9][0-9]{0,6}$`)
	d.AddValidator("userID", func(value string) error {
		if !userIDReg.MatchString(value) {
			return hes.New("user id should be numbers")
		}
		return nil
	})
}
