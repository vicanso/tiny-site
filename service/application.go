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
	"os"

	"go.uber.org/atomic"
)

var (
	// 应用状态
	applicationStatus = atomic.NewInt32(ApplicationStatusStopped)
	// 应用版本
	applicationVersion string
	// 应用构建时间
	applicationBuildedAt string
	// 应用运行环境的hostname
	applicationHostname string
)

const (
	// ApplicationStatusStopped 应用停止
	ApplicationStatusStopped int32 = iota + 1
	// ApplicationStatusRunning 应用运行中
	ApplicationStatusRunning
	// ApplicationStatusStopping 应用正在停止
	ApplicationStatusStopping
)

func init() {
	applicationHostname, _ = os.Hostname()
}

// SetApplicationStatus 设置应用运行状态
func SetApplicationStatus(status int32) {
	applicationStatus.Store(status)
}

// GetApplicationStatus 获取应用运行状态
func GetApplicationStatus() int32 {
	return applicationStatus.Load()
}

// ApplicationIsRunning 判断应用是否正在运行
func ApplicationIsRunning() bool {
	return applicationStatus.Load() == ApplicationStatusRunning
}

// GetApplicationVersion 获取应用版本号
func GetApplicationVersion() string {
	return applicationVersion
}

// SetApplicationVersion 设置应用版本号
func SetApplicationVersion(v string) {
	applicationVersion = v
}

// GetApplicationBuildedAt 获取应用构建时间
func GetApplicationBuildedAt() string {
	return applicationBuildedAt
}

// SetApplicationBuildedAt 设置应用构建时间
func SetApplicationBuildedAt(v string) {
	applicationBuildedAt = v
}

// GetApplicationHostname 获取应用运行的hostname
func GetApplicationHostname() string {
	return applicationHostname
}
