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
	"image"

	"github.com/disintegration/imaging"
	"github.com/vicanso/tiny-site/storage"
)

type resizeHandler func(image.Image, int, int, imaging.ResampleFilter) *image.NRGBA

func resize(fn resizeHandler, img *storage.Image, width, height int) (*storage.Image, error) {
	if img.Width <= width && img.Height <= height {
		return img, nil
	}
	srcImage, err := decodeImage(img)
	if err != nil {
		return nil, err
	}
	srcImage = fn(srcImage, width, height, imaging.Lanczos)
	data, err := encodeImage(srcImage, img.Type)
	if err != nil {
		return nil, err
	}
	img.Width = srcImage.Bounds().Dx()
	img.Height = srcImage.Bounds().Dy()
	img.SetData(data)
	return img, nil
}

func NewFitResizeImage(width, height int) ImageJob {
	return func(_ context.Context, img *storage.Image) (*storage.Image, error) {
		return resize(imaging.Fit, img, width, height)
	}
}

func NewFillResizeImage(width, height int) ImageJob {
	return func(_ context.Context, img *storage.Image) (*storage.Image, error) {
		return resize(func(i1 image.Image, i2, i3 int, rf imaging.ResampleFilter) *image.NRGBA {
			return imaging.Fill(i1, i2, i3, imaging.Center, rf)
		}, img, width, height)
	}
}
