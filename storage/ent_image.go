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
	"github.com/vicanso/tiny-site/ent/image"
	"github.com/vicanso/tiny-site/helper"
)

type entStorage struct {
	client *ent.Client
}

func mustNewEntStorage() ImageStorage {
	return &entStorage{
		client: helper.EntGetClient(),
	}
}

// Get gets image from ent(mysql or postgres)
func (e *entStorage) Get(ctx context.Context, bucket, name string) (*ent.Image, error) {
	return e.client.Image.Query().
		Where(image.BucketEQ(bucket)).
		Where(image.NameEQ(name)).
		First(ctx)
}

func (e *entStorage) update(ctx context.Context, data ent.Image) error {
	updateOne := e.client.Image.UpdateOneID(data.ID)
	if data.Bucket != "" {
		updateOne.SetBucket(data.Bucket)
	}
	if data.Name != "" {
		updateOne.SetName(data.Name)
	}
	if data.Type != "" {
		updateOne.SetType(data.Type)
	}
	if data.Width != 0 {
		updateOne.SetWidth(data.Width)
	}
	if data.Height != 0 {
		updateOne.SetHeight(data.Height)
	}
	if data.Metadata != nil {
		updateOne.SetMetadata(data.Metadata)
	}
	if len(data.Tags) != 0 {
		updateOne.SetTags(data.Tags)
	}
	size := len(data.Data)
	if size != 0 {
		updateOne.SetData(data.Data)
		updateOne.SetSize(size)
	}
	_, err := updateOne.Save(ctx)
	return err
}

// Put puts image to ent(mysql or postgres)
func (e *entStorage) Put(ctx context.Context, data ent.Image) error {
	// 如果指定了id，则是更新
	if data.ID != 0 {
		return e.update(ctx, data)
	}
	_, err := e.client.Image.Create().
		SetBucket(data.Bucket).
		SetName(data.Name).
		SetType(data.Type).
		SetSize(len(data.Data)).
		SetWidth(data.Width).
		SetHeight(data.Height).
		SetMetadata(data.Metadata).
		SetCreator(data.Creator).
		SetData(data.Data).
		SetTags(data.Tags).
		Save(ctx)
	return err
}

// Query gets the files from ent(mysql or postgres)
func (e *entStorage) Query(ctx context.Context, param ImageFilterParams) ([]*ent.Image, error) {
	return nil, nil
}

// Count counts the files from ent(mysql or postgres)
func (e *entStorage) Count(ctx context.Context, params ImageFilterParams) (int64, error) {
	return -1, nil
}
