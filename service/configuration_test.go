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

package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/tiny-site/cs"
)

func TestConfigurationService(t *testing.T) {
	t.Run("is valid", func(t *testing.T) {
		assert := assert.New(t)
		conf := &Configuration{
			Status: cs.ConfigDiabled,
		}
		assert.False(conf.IsValid())

		conf.Status = cs.ConfigEnabled
		d := time.Now()
		d = d.AddDate(0, 1, 0)
		oneMonthLater := d

		d = time.Now()
		d = d.AddDate(0, -1, 0)
		oneMonthAgo := d

		conf.BeginDate = &oneMonthLater
		assert.False(conf.IsValid())

		conf.BeginDate = &oneMonthAgo
		conf.EndDate = &oneMonthAgo
		assert.False(conf.IsValid())

		conf.EndDate = &oneMonthLater
		assert.True(conf.IsValid())
	})
}
