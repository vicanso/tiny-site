// Copyright 2020 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, softwareEntGetClient
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/hes"
)

func TestParseColor(t *testing.T) {
	assert := assert.New(t)
	c, err := parseColor("255,255,255")
	assert.Nil(err)
	assert.Equal(uint8(255), c.R)
	assert.Equal(uint8(255), c.G)
	assert.Equal(uint8(255), c.B)

	_, err = parseColor("255,255")
	assert.Equal("非法颜色值，格式必须为：1,1,1，当前为：255,255", err.(*hes.Error).Message)
	_, err = parseColor("255,255,256")
	assert.Equal("非法颜色值，必须>=0 <=255，当前为：256", err.(*hes.Error).Message)
}

func TestGetCaptcha(t *testing.T) {
	assert := assert.New(t)
	info, err := GetCaptcha(context.TODO(), "0,0,0", "255,255,255")
	assert.Nil(err)
	assert.Equal(4, len(info.Value))
	assert.GreaterOrEqual(time.Now().Add(5*time.Minute).Unix(), info.ExpiredAt.Unix())
	assert.NotEmpty(info.Data)
	assert.Equal("jpeg", info.Type)
}

func TestValidateCaptcha(t *testing.T) {
	assert := assert.New(t)
	info, err := GetCaptcha(context.TODO(), "0,0,0", "255,255,255")
	assert.Nil(err)
	valid, err := ValidateCaptcha(context.TODO(), info.ID, info.Value)
	assert.Nil(err)
	assert.True(valid)
}
