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

package schema

import (
	"net/http"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Image struct {
	ent.Schema
}

// Mixin 图片的mixin
func (Image) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

func (Image) Fields() []ent.Field {
	return []ent.Field{
		field.String("bucket").
			NotEmpty().
			Comment("图片所以bucket"),
		field.String("name").
			NotEmpty().
			Comment("图片名"),
		field.String("type").
			NotEmpty().
			Comment("图片类型"),
		field.Int("size").
			NonNegative().
			Comment("图片数据长度"),
		field.Int("width").
			NonNegative().
			Comment("图片宽度"),
		field.Int("height").
			NonNegative().
			Comment("图片高度"),
		field.String("tags").
			Comment("图片标签"),
		field.JSON("metadata", &http.Header{}).
			Optional().
			Comment("metadata"),
		field.String("creator").
			NotEmpty().
			Comment("创建者"),
		field.Bytes("data").
			Comment("图片数据"),
	}
}

// Indexes 文件表索引
func (Image) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("bucket", "name").Unique(),
	}
}
