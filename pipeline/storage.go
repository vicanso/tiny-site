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

	"github.com/vicanso/tiny-site/ent"
	entImage "github.com/vicanso/tiny-site/ent/image"
	"github.com/vicanso/tiny-site/helper"
)

func NewGetEntImage(bucket, name string) ImageJob {
	return func(ctx context.Context, _ *ent.Image) (*ent.Image, error) {
		return helper.EntGetClient().Image.Query().
			Where(entImage.Bucket(bucket)).
			Where(entImage.Name(name)).
			First(ctx)
	}
}
