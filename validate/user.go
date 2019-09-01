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
	"github.com/asaskevich/govalidator"

	"github.com/vicanso/tiny-site/cs"
)

func init() {
	// 账号
	Add("xUserAccount", func(i interface{}, _ interface{}) bool {
		return checkASCIIStringLength(i, 4, 10)
	})
	Add("xUserPassword", func(i interface{}, _ interface{}) bool {
		return checkASCIIStringLength(i, 44, 44)
	})
	Add("xUserAccountKeyword", func(i interface{}, _ interface{}) bool {
		return checkASCIIStringLength(i, 1, 10)
	})
	Add("xUserRole", func(i interface{}, _ interface{}) bool {
		value, ok := i.(string)
		if !ok {
			return false
		}
		return govalidator.IsIn(value, cs.UserRoleSu, cs.UserRoleAdmin)
	})
	Add("xUserRoles", func(i interface{}, _ interface{}) bool {
		values, ok := i.([]string)
		if !ok {
			return false
		}
		for _, value := range values {
			if !govalidator.IsIn(value, cs.UserRoleSu, cs.UserRoleAdmin) {
				return false
			}
		}
		return true
	})
}
