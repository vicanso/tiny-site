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

package validate

import (
	"github.com/vicanso/tiny-site/cs"
)

func init() {
	// 应用配置名称
	Add("xConfigName", func(i interface{}, _ interface{}) bool {
		return checkAlphanumericStringLength(i, 2, 20)
	})
	Add("xConfigCategory", func(i interface{}, _ interface{}) bool {
		return checkAlphanumericStringLength(i, 2, 20)
	})
	Add("xConfigData", func(i interface{}, _ interface{}) bool {
		return checkStringLength(i, 1, 500)
	})
	Add("xConfigNames", func(i interface{}, _ interface{}) bool {
		return checkAlphanumericStringLength(i, 2, 100)
	})
	Add("xConfigStatus", func(i interface{}, _ interface{}) bool {
		value, ok := i.(int)
		if !ok {
			return false
		}
		return value == cs.ConfigEnabled || value == cs.ConfigDiabled
	})
}
