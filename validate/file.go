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
)

func init() {
	Add("xFileZoneName", func(i interface{}, _ interface{}) bool {
		return checkStringLength(i, 4, 20)
	})

	Add("xFileZone", func(i interface{}, _ interface{}) bool {
		value := govalidator.ToString(i)
		return govalidator.Range(value, "1", "1000")
	})

	Add("xFileZoneAuthority", func(i interface{}, _ interface{}) bool {
		value := govalidator.ToString(i)
		return govalidator.IsIn(value, "1", "2")
	})
	Add("xFileZoneDesc", func(i interface{}, _ interface{}) bool {
		return checkStringLength(i, 1, 100)
	})

	Add("xFileName", func(i interface{}, _ interface{}) bool {
		return checkStringLength(i, 4, 26)
	})

	Add("xFileType", func(i interface{}, _ interface{}) bool {
		return checkASCIIStringLength(i, 1, 5)
	})
	Add("xFileDesc", func(i interface{}, _ interface{}) bool {
		return checkStringLength(i, 1, 100)
	})
}
