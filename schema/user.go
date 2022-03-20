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
	"regexp"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// 用户角色
const (
	// UserRoleNormal normal user
	UserRoleNormal = "normal"
	// UserRoleSu super user
	UserRoleSu = "su"
	// UserRoleAdmin admin user
	UserRoleAdmin = "admin"
)

// UserRoleInfo 用户角色信息
type UserRoleInfo struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// GetUserRoleList 获取用户角色列表
func GetUserRoleList() []*UserRoleInfo {
	return []*UserRoleInfo{
		{
			Name:  "普通用户",
			Value: UserRoleNormal,
		},
		{
			Name:  "管理员用户",
			Value: UserRoleAdmin,
		},
		{
			Name:  "超级用户",
			Value: UserRoleSu,
		},
	}
}

// Mixin 用户表的mixin
func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
		StatusMixin{},
	}
}

// Fields 用户表的字段配置
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("account").
			Match(regexp.MustCompile("[a-zA-Z_0-9]+$")).
			NotEmpty().
			Immutable().
			Unique().
			Comment("用户账户信息"),
		field.String("password").
			Sensitive().
			NotEmpty().
			Comment("用户密码，保存hash之后的值"),
		field.String("name").
			Optional().
			Comment("用户名称"),
		field.Strings("roles").
			Optional().
			Comment("用户角色，由管理员分配"),
		field.Strings("groups").
			Optional().
			Comment("用户分组，按用户职能分配至不同的分组"),
		field.String("email").
			Optional().
			Comment("用户邮箱"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}

// Indexes 用户表索引
func (User) Indexes() []ent.Index {
	return []ent.Index{
		// 用户账户唯一索引
		index.Fields("account").Unique(),
	}
}
