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

package pipeline

import (
	"context"

	entImage "github.com/vicanso/tiny-site/ent/image"
	"github.com/vicanso/tiny-site/helper"
	"github.com/vicanso/tiny-site/storage"
)

func NewGetEntImage(bucket, name string) ImageJob {
	return func(ctx context.Context, _ *storage.Image) (*storage.Image, error) {
		img, err := helper.EntGetClient().Image.Query().
			Where(entImage.Bucket(bucket)).
			Where(entImage.Name(name)).
			First(ctx)
		if err != nil {
			return nil, err
		}
		return &storage.Image{
			Type:   img.Type,
			Size:   img.Size,
			Width:  img.Width,
			Height: img.Height,
			Data:   img.Data,
		}, nil
	}
}

func NewGetHTTPImage(url string) ImageJob {
	return func(ctx context.Context, _ *storage.Image) (*storage.Image, error) {
		return storage.GetImageFromURL(ctx, url)
	}
}
