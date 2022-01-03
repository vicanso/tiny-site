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

package storage

import (
	"context"

	"github.com/vicanso/tiny-site/ent"
)

type ImageFilterParams struct {
	// 筛选的字段
	Fields string `json:"fields"`
	// 数量
	Limit int `json:"limit"`
	// 偏移量
	Offset int `json:"offset"`
}

type ImageStorage interface {
	Get(ctx context.Context, bucket, filename string) (*ent.Image, error)
	Put(ctx context.Context, file ent.Image) error
	Query(ctx context.Context, params ImageFilterParams) ([]*ent.Image, error)
	Count(ctx context.Context, params ImageFilterParams) (int64, error)
}

var entStorageClient = mustNewEntStorage()

func Ent() ImageStorage {
	return entStorageClient
}