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

package service

import (
	"github.com/vicanso/ips"
)

var blockIPS = ips.New()

// ResetIPBlocker 重置IP拦截器的IP列表
func ResetIPBlocker(ipList []string) error {
	return blockIPS.Replace(ipList...)
}

// IsBlockIP 判断该IP是否有需要拦截
func IsBlockIP(ip string) bool {
	return blockIPS.Contains(ip)
}

// GetIPBlockList 获取block的ip地址列表
func GetIPBlockList() []string {
	return blockIPS.Strings()
}
