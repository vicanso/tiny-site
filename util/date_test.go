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

package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNow(t *testing.T) {
	assert := assert.New(t)
	defer SetMockTime("")
	mockTime := "2020-04-26T20:34:33+08:00"
	SetMockTime(mockTime)

	assert.Equal(int64(1587904473000000000), Now().UnixNano())
	// travis中为0时区
	// assert.Equal(mockTime, NowString())

	assert.Equal("2020-04-26 12:34:33 +0000 UTC", UTCNow().String())
	assert.Equal("2020-04-26T12:34:33Z", UTCNowString())

	value, err := ParseTime(mockTime)
	assert.Nil(err)
	assert.Equal("2020-04-26T20:34:33+08:00", FormatTime(value))

	chinaNow, err := ChinaNow()
	assert.Nil(err)
	assert.Equal("2020-04-26T20:34:33+08:00", FormatTime(chinaNow))

	chinaToday, err := ChinaToday()
	assert.Nil(err)
	assert.Equal("2020-04-26T00:00:00+08:00", FormatTime(chinaToday))

	chinaYesterday, err := ChinaYesterday()
	assert.Nil(err)
	assert.Equal("2020-04-25T00:00:00+08:00", FormatTime(chinaYesterday))
}

func TestIsBetween(t *testing.T) {
	assert := assert.New(t)
	defer SetMockTime("")
	mockTime := "2020-04-26T20:34:33+08:00"
	SetMockTime(mockTime)

	start, _ := ParseTime("2020-04-26T19:34:33+08:00")
	end, _ := ParseTime("2020-04-26T21:34:33+08:00")
	assert.True(IsBetween(start, end))
	assert.False(IsBetween(start, start))
	assert.False(IsBetween(end, end))
}

func TestNewTimeWithRandomNS(t *testing.T) {
	assert := assert.New(t)
	date := NewTimeWithRandomNS(0)
	assert.Equal(int64(0), date.Unix())
}
