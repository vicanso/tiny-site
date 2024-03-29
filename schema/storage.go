// Copyright 2022 tree xie
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
)

const (
	StorageCategoryHTTP   = "http"
	StorageCategoryMinio  = "minio"
	StorageCategoryOSS    = "oss"
	StorageCategoryGridfs = "gridfs"
)

type Storage struct {
	ent.Schema
}

func (Storage) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
		StatusMixin{},
	}
}

func (Storage) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Unique().
			Comment("存储服务名称"),
		field.Enum("category").
			Values(
				StorageCategoryHTTP,
				StorageCategoryMinio,
				StorageCategoryOSS,
				StorageCategoryGridfs,
			).
			Comment("存储类型"),
		field.Text("uri").
			NotEmpty().
			Comment("存储连接串"),
		field.Text("description").
			Optional().
			Comment("存储描述"),
	}
}
