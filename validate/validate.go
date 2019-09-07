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
	"encoding/json"
	"regexp"

	"github.com/asaskevich/govalidator"
	jsoniter "github.com/json-iterator/go"
	"github.com/vicanso/hes"
)

var (
	standardJSON         = jsoniter.ConfigCompatibleWithStandardLibrary
	paramTagRegexMap     = govalidator.ParamTagRegexMap
	paramTagMap          = govalidator.ParamTagMap
	customTypeTagMap     = govalidator.CustomTypeTagMap
	errCategory          = "validate"
	errJSONParseCategory = "json-parse"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func doValidate(s interface{}, data interface{}) (err error) {
	// statusCode := http.StatusBadRequest
	if data != nil {
		switch data.(type) {
		case []byte:
			err = json.Unmarshal(data.([]byte), s)
			if err != nil {
				he := hes.Wrap(err)
				he.Category = errJSONParseCategory
				err = he
				return
			}
		default:
			buf, err := json.Marshal(data)
			if err != nil {
				return err
			}
			err = json.Unmarshal(buf, s)
			if err != nil {
				return err
			}
		}
	}
	_, err = govalidator.ValidateStruct(s)
	return
}

// Do do validate
func Do(s interface{}, data interface{}) (err error) {
	err = doValidate(s, data)
	if err != nil {
		he := hes.Wrap(err)
		if he.Category == "" {
			he.Category = errCategory
		}
		err = he
	}
	return
}

// AddRegex add a regexp validate
func AddRegex(name, reg string, fn govalidator.ParamValidator) {
	if paramTagMap[name] != nil {
		panic(name + ", reg:" + reg + " is duplicated")
	}
	paramTagRegexMap[name] = regexp.MustCompile(reg)
	paramTagMap[name] = fn
}

// Add add validate
func Add(name string, fn govalidator.CustomTypeValidator) {
	_, exists := customTypeTagMap.Get(name)
	if exists {
		panic(name + " is duplicated")
	}
	customTypeTagMap.Set(name, fn)
}

func checkASCIIStringLength(i interface{}, min, max int) bool {
	value, ok := i.(string)
	if !ok {
		return false
	}
	if !govalidator.IsASCII(value) {
		return false
	}
	size := len(value)
	if size < min || size > max {
		return false
	}
	return true
}

func checkAlphaStringLength(i interface{}, min, max int) bool {
	value, ok := i.(string)
	if !ok {
		return false
	}
	if !govalidator.IsAlpha(value) {
		return false
	}
	size := len(value)
	if size < min || size > max {
		return false
	}
	return true
}

func checkAlphanumericStringLength(i interface{}, min, max int) bool {
	value, ok := i.(string)
	if !ok {
		return false
	}
	if !govalidator.IsAlphanumeric(value) {
		return false
	}
	size := len(value)
	if size < min || size > max {
		return false
	}
	return true
}

func checkStringLength(i interface{}, min, max int) bool {
	value, ok := i.(string)
	if !ok {
		return false
	}
	size := len(value)
	if size < min || size > max {
		return false
	}
	return true
}
