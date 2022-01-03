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
	return func(_ context.Context, i *ent.Image) (*ent.Image, error) {
		if i.Width <= width && i.Height <= height {
			return i, nil
		}
		img, _, err := image.Decode(bytes.NewReader(i.Data))
		if err != nil {
			return nil, err
		}
		img = imaging.Fit(img, width, height, imaging.Lanczos)
		buffer := bytes.Buffer{}
		format := imaging.JPEG
		if i.Type == "png" {
			format = imaging.PNG
		}
		err = imaging.Encode(&buffer, img, format)
		if err != nil {
			return nil, err
		}
		i.Width = img.Bounds().Dx()
		i.Height = img.Bounds().Dy()
		i.Data = buffer.Bytes()
		return i, nil
	}
}
