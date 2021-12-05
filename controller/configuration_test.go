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

package controller

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/tiny-site/ent/configuration"
	confSchema "github.com/vicanso/tiny-site/ent/configuration"
	"github.com/vicanso/tiny-site/helper"
	"github.com/vicanso/tiny-site/schema"
	"github.com/vicanso/tiny-site/util"
	"github.com/vicanso/hes"
)

func TestConfigurationParams(t *testing.T) {
	// 测试代码中不执行main，因此在此处调用初始化
	_ = helper.EntInitSchema()
	assert := assert.New(t)
	now := func() time.Time {
		t := time.Now()
		return t
	}

	name := util.RandomString(8)
	category := confSchema.CategoryBlockIP
	defer func() {
		_, _ = getConfigurationClient().Delete().Where(configuration.Name(name)).Exec(context.Background())
	}()

	configID := 0

	t.Run("save", func(t *testing.T) {
		params := configurationAddParams{
			Name:      name,
			Status:    schema.StatusDisabled,
			Category:  category,
			StartedAt: now(),
			EndedAt:   now(),
			Data:      "test",
		}
		conf, err := params.save(context.Background(), "treexie")
		assert.Nil(err)
		assert.Equal(name, conf.Name)
		configID = conf.ID
	})

	t.Run("validateBeforeSave", func(t *testing.T) {
		params := configurationAddParams{
			Name: util.RandomString(8),
		}
		err := params.validateBeforeSave(context.Background())
		assert.Nil(err)
		params.Name = name
		err = params.validateBeforeSave(context.Background())
		assert.Equal("该配置已存在", err.(*hes.Error).Message)
	})

	t.Run("query by name", func(t *testing.T) {
		params := configurationListParmas{
			Name: name,
		}
		configs, err := params.queryAll(context.Background())
		assert.Nil(err)
		assert.Equal(1, len(configs))
	})

	t.Run("query by category", func(t *testing.T) {
		params := configurationListParmas{
			Category: category,
		}
		params.listParams.Order = "-created_at"
		params.listParams.Limit = 1
		configs, err := params.queryAll(context.Background())
		assert.Nil(err)
		assert.Equal(1, len(configs))
		assert.Equal(name, configs[0].Name)
	})

	t.Run("count", func(t *testing.T) {
		params := configurationListParmas{}
		count, err := params.count(context.Background())
		assert.Nil(err)
		assert.True(count >= 1)
	})

	t.Run("update", func(t *testing.T) {
		newTime := now()
		status := schema.StatusEnabled
		category := confSchema.CategoryMockTime
		data := "mock time"
		params := configurationUpdateParams{
			StartedAt: newTime,
			EndedAt:   newTime,
			Status:    status,
			Category:  category,
			Data:      data,
		}
		conf, err := params.updateOneID(context.Background(), configID)
		assert.Nil(err)
		assert.Equal(status, conf.Status)
		assert.Equal(category, conf.Category)
		assert.Equal(data, conf.Data)
		assert.Equal(newTime.Unix(), conf.StartedAt.Unix())
		assert.Equal(newTime.Unix(), conf.EndedAt.Unix())
	})
}
