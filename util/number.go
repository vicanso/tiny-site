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

package util

import (
	"fmt"
	"math"
	"strconv"
)

// 可根据需求调整相差值的判断
const tolerance = 1e-6

// AlmostEqual returns true is abs(a - b) < 1e-6
func AlmostEqual(a, b float64) bool {
	return math.Abs(a-b) <= tolerance
}

// ToFixed formats float64 to string
func ToFixed(value float64, precision int) string {
	str := "%." + strconv.Itoa(precision) + "f"
	return fmt.Sprintf(str, value)
}
