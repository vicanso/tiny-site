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
	"bytes"
	"time"

	"github.com/felixge/fgprof"
	"github.com/vicanso/hes"
)

func GetProf(d time.Duration) (*bytes.Buffer, error) {
	// 禁止拉取超过1分钟的prof
	if d > 1*time.Minute {
		return nil, hes.New("duration should be less than 1m")
	}
	result := &bytes.Buffer{}
	done := fgprof.Start(result, fgprof.FormatPprof)
	time.Sleep(d)
	err := done()
	if err != nil {
		return nil, err
	}
	return result, nil
}
