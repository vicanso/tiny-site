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
	"math/rand"
	"time"

	"github.com/jinzhu/now"
	"go.uber.org/atomic"
)

var mockTime atomic.Int64

// SetMockTime 设置mock的时间
func SetMockTime(v string) {
	if v == "" {
		mockTime.Store(0)
		return
	}
	t, _ := time.Parse(time.RFC3339, v)
	seconds := t.Unix()
	// 设置的时间有误，不调整
	if seconds < 0 {
		return
	}
	mockTime.Store(seconds)
}

// GetMockTime 获取mock的时间
func GetMockTime() string {
	v := mockTime.Load()
	if v == 0 {
		return ""
	}
	return FormatTime(time.Unix(v, 0))
}

// Now 获取当前时间（测试环境允许使用mock的时间)
func Now() time.Time {
	// 正式环境不提供mock time
	if IsProduction() {
		return time.Now()
	}
	v := mockTime.Load()
	if v == 0 {
		return time.Now()
	}
	return time.Unix(v, 0)
}

// NowString 获取当前时间字符串RFC3339
func NowString() string {
	return Now().Format(time.RFC3339)
}

// UTCNow 获取UTC的时间
func UTCNow() time.Time {
	return Now().UTC()
}

// UTCNowString 获取UTC时间字符串RFC3339
func UTCNowString() string {
	return UTCNow().Format(time.RFC3339)
}

// ParseTime 字符串转换为时间，字符串为RFC3339格式
func ParseTime(str string) (time.Time, error) {
	return time.Parse(time.RFC3339, str)
}

// FormatTime 格式化时间为字符串RFC3339
func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ChinaNow 获取中国时间
func ChinaNow() (time.Time, error) {
	t := Now()
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return t, err
	}
	return t.In(loc), nil
}

// ChinaToday 获取中国时间当天的起始时间
func ChinaToday() (time.Time, error) {
	t, err := ChinaNow()
	return now.With(t).BeginningOfDay(), err
}

// ChinaYesterday 获取中国时间昨天的起始时间
func ChinaYesterday() (time.Time, error) {
	t, err := ChinaNow()
	if err != nil {
		return t, err
	}
	t = t.AddDate(0, 0, -1)
	return now.With(t).BeginningOfDay(), err
}

// IsBetween 判断当前时间是否在开始与结束时间之间
func IsBetween(begin, end time.Time) bool {
	now := Now().Unix()
	// 如果开始时间大于当前时间，未开始
	if !begin.IsZero() && begin.Unix() > now {
		return false
	}
	// 如果结束时间少于当前时间，已结束
	if !end.IsZero() && end.Unix() < now {
		return false
	}
	return true
}

// NewTimeWithRandomNS 根据timestamp并添加随机的ns生成时间
func NewTimeWithRandomNS(timestamp int64) time.Time {
	rand.Seed(time.Now().UnixNano())
	sec := timestamp / 1000
	ms := timestamp % 1000
	ns := ms*10e6 + time.Now().UnixNano()%10e6
	return time.Unix(sec, ns)
}
