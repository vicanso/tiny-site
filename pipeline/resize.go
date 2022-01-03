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
	"bytes"
	"context"
	"image"

	"github.com/disintegration/imaging"
	"github.com/vicanso/tiny-site/ent"
)

func NewFitResizeImage(width, height int) ImageJob {
	return func(_ context.Context, img *ent.Image) (*ent.Image, error) {
		if img.Width <= width && img.Height <= height {
			return img, nil
		}
		srcImage, _, err := image.Decode(bytes.NewReader(img.Data))
		if err != nil {
			return nil, err
		}
		srcImage = imaging.Fit(srcImage, width, height, imaging.Lanczos)
		buffer := bytes.Buffer{}
		format := imaging.JPEG
		if img.Type == "png" {
			format = imaging.PNG
		}
		err = imaging.Encode(&buffer, srcImage, format)
		if err != nil {
			return nil, err
		}
		img.Width = srcImage.Bounds().Dx()
		img.Height = srcImage.Bounds().Dy()
		img.Data = buffer.Bytes()
		return img, nil
	}
}
