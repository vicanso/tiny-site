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

const (
	// ConfigurationCategoryMockTime mock time分类
	ConfigurationCategoryMockTime = "mockTime"
	// ConfigurationCategoryBlockIP block ip分类
	ConfigurationCategoryBlockIP = "blockIP"
	// ConfigurationCategorySignedKey signed key分类
	ConfigurationCategorySignedKey = "signedKey"
	// ConfigurationCategoryRouterConcurrency router concurrency分类
	ConfigurationCategoryRouterConcurrency = "routerConcurrency"
	// ConfigurationCategoryRouter router分类
	ConfigurationCategoryRouter = "router"
	// ConfigurationCategoryRequestConcurrency request concurrency
	ConfigurationCategoryRequestConcurrency = "requestConcurrency"
	// ConfigurationCategoryEmail 邮箱配置
	ConfigurationCategoryEmail = "email"
	// ConfigurationHTTPServerInterceptor http服务的拦截配置
	ConfigurationHTTPServerInterceptor = "httpServerInterceptor"
)

// Configuration holds the schema definition for the Configuration entity.
type Configuration struct {
	ent.Schema
}

// Mixin 配置信息的mixin
func (Configuration) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
		StatusMixin{},
	}
}

// Fields 配置信息的相关字段
func (Configuration) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Unique().
			Comment("配置名称"),
		field.Enum("category").
			Values(
				ConfigurationCategoryMockTime,
				ConfigurationCategoryBlockIP,
				ConfigurationCategorySignedKey,
				ConfigurationCategoryRouterConcurrency,
				ConfigurationCategoryRouter,
				ConfigurationCategoryRequestConcurrency,
				ConfigurationCategoryEmail,
				ConfigurationHTTPServerInterceptor,
			).
			Comment("配置分类"),
		field.String("owner").
			NotEmpty().
			Comment("创建者"),
		field.Text("data").
			NotEmpty().
			Comment("配置信息"),
		field.Time("started_at").
			StructTag(`json:"startedAt"`).
			Comment("配置启用时间"),
		field.Time("ended_at").
			StructTag(`json:"endedAt"`).
			Comment("配置停用时间"),
		field.String("description").
			Comment("配置说明").
			Optional(),
	}
}

// Edges of the Configuration.
func (Configuration) Edges() []ent.Edge {
	return nil
}

// Indexes 配置表索引
func (Configuration) Indexes() []ent.Index {
	return []ent.Index{
		// 配置名称，唯一索引
		index.Fields("name").Unique(),
		index.Fields("status"),
	}
}
