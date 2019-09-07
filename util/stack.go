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
	"regexp"
	"runtime"
	"strings"
)

var (
	stackReg = regexp.MustCompile(`\((0x[\s\S]+)\)`)
	appPath  = ""
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	fileDivide := "/"
	arr := strings.Split(file, fileDivide)
	arr = arr[0 : len(arr)-2]
	appPath = strings.Join(arr, fileDivide) + fileDivide
}

func getStack(length int, ignoreMod bool) []string {
	size := 2 << 10
	stack := make([]byte, size)
	runtime.Stack(stack, true)
	arr := strings.Split(string(stack), "\n")
	arr = arr[1:]
	max := len(arr) - 1
	result := []string{}
	for index := 0; index < max; index += 2 {
		if index+1 >= max {
			break
		}
		if ignoreMod && !strings.Contains(arr[index+1], appPath) {
			continue
		}

		file := strings.Replace(arr[index+1], appPath, "", 1)
		tmpArr := strings.Split(arr[index], "/")
		fn := stackReg.ReplaceAllString(tmpArr[len(tmpArr)-1], "")
		str := fn + ": " + file
		result = append(result, str)
	}
	if len(result) < length {
		return result
	}
	return result[:length]
}

// GetStack 获取调用信息
func GetStack(max int) []string {
	// 去除get stack的两个函数调用
	arr := getStack(max+2, true)
	return arr[2:]
}
