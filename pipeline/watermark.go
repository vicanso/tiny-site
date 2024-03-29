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
	"image/color"

	"github.com/disintegration/imaging"
	"github.com/vicanso/tiny-site/storage"
)

func NewWatermark(url string, postion string, angle float64) ImageJob {
	return func(ctx context.Context, img *storage.Image) (*storage.Image, error) {
		info, err := storage.GetImageFromURL(ctx, url)
		if err != nil {
			return nil, err
		}
		watermarkImg, err := info.Image()
		if err != nil {
			return nil, err
		}
		if angle != 0 {
			watermarkImg = imaging.Rotate(watermarkImg, angle, color.Transparent)
		}
		dst, err := decodeImage(img)
		if err != nil {
			return nil, err
		}
		x := 0
		y := 0
		watermarkWidth := watermarkImg.Bounds().Dx()
		watermarkHeight := watermarkImg.Bounds().Dy()
		switch postion {
		case PositionTop:
			x = (img.Width - watermarkWidth) / 2
		case PositionTopRight:
			x = img.Width - watermarkWidth
		case PositionLeft:
			y = (img.Height - watermarkHeight) / 2
		case PositionCenter:
			x = (img.Width - watermarkWidth) / 2
			y = (img.Height - watermarkHeight) / 2
		case PositionRight:
			y = (img.Height - watermarkHeight) / 2
			x = img.Width - watermarkWidth
		case PositionBottomLeft:
			y = img.Height - watermarkHeight
		case PositionBottom:
			x = (img.Width - watermarkWidth) / 2
			y = img.Height - watermarkHeight
		case PositionBottomRight:
			x = img.Width - watermarkWidth
			y = img.Height - watermarkHeight
		}
		dst = imaging.Paste(dst, watermarkImg, image.Pt(x, y))
		data, err := encodeImage(dst, img.Type)
		if err != nil {
			return nil, err
		}
		img.SetData(data)

		return img, nil
	}
}
