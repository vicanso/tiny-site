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
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Bucket struct {
	ent.Schema
}

func (Bucket) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

func (Bucket) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Immutable().
			Comment("bucket的名称"),
		field.String("creator").
			NotEmpty().
			Comment("bucket的创建者"),
		// 为空表示所有人可使用
		// 不为空则该列表中的用户可使用
		field.Strings("owners").
			Optional().
			Comment("bucket的拥有者"),
		field.String("description").
			NotEmpty().
			Comment("bucket的描述"),
	}
}

func (Bucket) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
	}
}
