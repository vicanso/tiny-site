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
	"bytes"
	"context"
	"image"

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

type ImageFinder func(ctx context.Context, params ...string) (*Image, error)

type Image struct {
	// 图片类型
	Type string
	// 原始数据长度(字节)
	OriginalSize int
	// 处理后数据长度(字节)
	Size int
	// 图片宽度
	Width int
	// 图片高度
	Height int
	// 图片数据
	Data []byte
	// 图片数据转换的图像
	img image.Image
}

func (i *Image) Image() (image.Image, error) {
	if i.img == nil {
		img, _, err := image.Decode(bytes.NewReader(i.Data))
		if err != nil {
			return nil, err
		}
		i.img = img
	}
	return i.img, nil
}

func NewImageFromBytes(data []byte) (*Image, error) {
	img, t, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	size := len(data)
	return &Image{
		Type:         t,
		Size:         size,
		OriginalSize: size,
		Width:        img.Bounds().Dx(),
		Height:       img.Bounds().Dy(),
		Data:         data,
		img:          img,
	}, nil
}

func (i *Image) SetData(data []byte) {
	i.Data = data
	i.Size = len(data)
	i.img = nil
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
