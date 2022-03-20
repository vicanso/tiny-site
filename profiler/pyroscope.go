// +build pyroscope

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
package profiler

import (
	pyroscopeProfiler "github.com/pyroscope-io/pyroscope/pkg/agent/profiler"
	"github.com/vicanso/tiny-site/config"
)

// MustStartPyroscope 启动pyroscope
func MustStartPyroscope() {
	basicConfig := config.MustGetBasicConfig()
	pyroscopeConfig := config.MustGetPyroscopeConfig()
	pyroscopeProfiler.Start(pyroscopeProfiler.Config{
		ApplicationName: basicConfig.Name,
		ServerAddress:   pyroscopeConfig.Addr,
		AuthToken:       pyroscopeConfig.Token,
	})
}
