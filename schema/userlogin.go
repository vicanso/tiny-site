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
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// UserLogin holds the schema definition for the UserLogin entity.
type UserLogin struct {
	ent.Schema
}

// Mixin 用户登录记录的mixin
func (UserLogin) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields 用户登录表的相关字段
func (UserLogin) Fields() []ent.Field {
	return []ent.Field{
		field.String("account").
			NotEmpty().
			Immutable().
			Comment("登录账户"),
		field.String("user_agent").
			StructTag(`json:"userAgent"`).
			Optional().
			Comment("用户浏览器的user-agent"),
		field.String("ip").
			Optional().
			Comment("用户IP"),
		field.String("track_id").
			StructTag(`json:"trackID"`).
			Optional().
			Comment("用户的track id"),
		field.String("session_id").
			StructTag(`json:"sessionID"`).
			Optional().
			Comment("用户的session id"),
		field.String("x_forwarded_for").
			StructTag(`json:"xForwardedFor"`).
			Optional().
			Comment("用户登录时的x-forwarded-for"),
		field.String("country").
			Optional().
			Comment("用户登录IP定位的国家"),
		field.String("province").
			Optional().
			Comment("用户登录IP定位的省份"),
		field.String("city").
			Optional().
			Comment("用户登录IP定位的城市"),
		field.String("isp").
			Optional().
			Comment("用户登录IP的网络服务商"),
	}
}

// Edges of the UserLogin.
func (UserLogin) Edges() []ent.Edge {
	return nil
}

// Indexes 用户登录表索引
func (UserLogin) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("account"),
	}
}
