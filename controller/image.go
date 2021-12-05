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

package controller

import (
	"github.com/vicanso/elton"
	"github.com/vicanso/tiny-site/cs"
	"github.com/vicanso/tiny-site/router"
)

type imageCtrl struct{}

type (
	addBucketParams struct {
		// bucket的名称
		Name string `json:"name"`
		// 拥有者
		Owners []string `json:"owners"`
		// 描述
		Description string `json:"description"`
	}
)

func init() {
	g := router.NewGroup("/images", loadUserSession)
	ctrl := imageCtrl{}

	g.POST(
		"/v1/buckets",
		shouldBeLogin,
		newTrackerMiddleware(cs.ActionBucketCreate),
		ctrl.createBucket,
	)
}

func (*imageCtrl) createBucket(c *elton.Context) error {
	params := addBucketParams{}
	err := validateBody(c, &params)
	if err != nil {
		return err
	}
	bucket, err := getBucketClient().Create().
		SetName(params.Name).
		SetOwners(params.Owners).
		SetDescription(params.Description).
		Save(c.Context())
	if err != nil {
		return err
	}
	c.Created(bucket)
	return nil
}
