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
	"errors"

	"github.com/vicanso/tiny-site/ent"
)

// 不再执行后续时返回
var ErrAbort = errors.New("abort")

type ImageJob func(context.Context, *ent.Image) (*ent.Image, error)

func Do(ctx context.Context, img *ent.Image, jobs ...ImageJob) (*ent.Image, error) {
	var err error
	for _, fn := range jobs {
		img, err = fn(ctx, img)
		if err != nil {
			// 如果是abort error，则直接返回数据
			if err == ErrAbort {
				return img, nil
			}
			return nil, err
		}
	}
	return img, nil
}
