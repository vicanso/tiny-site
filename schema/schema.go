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

package schema

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/tidwall/gjson"
)

type Status int8

// StatusInfo 状态信息
type StatusInfo struct {
	Name  string `json:"name"`
	Value Status `json:"value"`
}

func (s *Status) MarshalJSON() ([]byte, error) {
	str := s.String()
	return json.Marshal(map[string]interface{}{
		"value": int(*s),
		"desc":  str,
	})
}

func (s *Status) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	var value int64
	// 兼容两种方式
	if regexp.MustCompile(`^\d+$`).Match(data) {
		v, _ := strconv.Atoi(string(data))
		value = int64(v)
	} else {
		value = gjson.GetBytes(data, "value").Int()
	}
	reflect.ValueOf(s).Elem().SetInt(value)
	return nil
}

const (
	// 状态启用
	StatusEnabled Status = iota + 1
	// 状态禁用
	StatusDisabled
)

// ToInt8 转换为int8
func (status Status) Int8() int8 {
	return int8(status)
}

// String 转换为string
func (status Status) String() string {
	switch status {
	case StatusEnabled:
		return "启用"
	case StatusDisabled:
		return "禁用"
	default:
		return "未知"
	}
}

// GetSchemaStatusList 获取schema的状态列表
func GetStatusList() []*StatusInfo {
	values := []Status{
		StatusEnabled,
		StatusDisabled,
	}
	list := make([]*StatusInfo, len(values))
	for index, value := range values {
		list[index] = &StatusInfo{
			Name:  value.String(),
			Value: value,
		}
	}
	return list
}

// TimeMixin 公共的时间schema
type TimeMixin struct {
	mixin.Schema
}

// Fields 公共时间schema的字段，包括创建于与更新于
func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			// 对于多个单词组成的，如果需要使用select，则需要添加sql tag
			StructTag(`json:"createdAt" sql:"created_at"`).
			Immutable().
			Default(time.Now).
			Comment("创建时间，添加记录时由程序自动生成"),
		field.Time("updated_at").
			StructTag(`json:"updatedAt" sql:"updated_at"`).
			Default(time.Now).
			Immutable().
			UpdateDefault(time.Now).
			Comment("更新时间，更新记录时由程序自动生成"),
	}
}

// Indexes 公共时间字段索引
func (TimeMixin) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("created_at"),
		index.Fields("updated_at"),
	}
}

// StatusMixin 状态的schema
type StatusMixin struct {
	mixin.Schema
}

// Fields 公共的status的字段
func (StatusMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Int8("status").
			Range(StatusEnabled.Int8(), StatusDisabled.Int8()).
			Default(StatusEnabled.Int8()).
			GoType(Status(StatusEnabled)).
			Comment("状态，默认为启用状态"),
	}
}
