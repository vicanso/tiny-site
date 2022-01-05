// Copyright 2022 tree xie
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
	"context"

	"github.com/vicanso/elton"
	"github.com/vicanso/tiny-site/cs"
	"github.com/vicanso/tiny-site/ent"
	"github.com/vicanso/tiny-site/ent/storage"
	"github.com/vicanso/tiny-site/router"
)

type storageCtrl struct{}

type (
	storageAddParams struct {
		Name     string `json:"name" validate:"required,xStorageName"`
		Category string `json:"category" validate:"required,xStorageCategory"`
		URI      string `json:"uri" validate:"required,xStorageURI"`
	}
)

type (
	storageListResp struct {
		Storages []*ent.Storage `json:"storages"`
	}
)

func init() {
	prefix := "/storages"

	g := router.NewGroup(prefix, loadUserSession)
	ctrl := storageCtrl{}

	g.GET(
		"/v1",
		shouldBeAdmin,
		ctrl.list,
	)

	g.POST(
		"/v1",
		newTrackerMiddleware(cs.ActionStorageAdd),
		shouldBeAdmin,
		ctrl.add,
	)
}

func (params *storageAddParams) save(ctx context.Context) (*ent.Storage, error) {
	return getStorageClient().Create().
		SetName(params.Name).
		SetCategory(storage.Category(params.Category)).
		SetURI(params.URI).
		Save(ctx)
}

func (*storageCtrl) add(c *elton.Context) error {
	params := storageAddParams{}
	err := validateBody(c, &params)
	if err != nil {
		return err
	}

	result, err := params.save(c.Context())
	if err != nil {
		return err
	}
	c.Created(result)
	return nil
}

func (*storageCtrl) list(c *elton.Context) error {
	storages, err := getStorageClient().Query().
		All(c.Context())

	if err != nil {
		return err
	}
	c.Body = &storageListResp{
		Storages: storages,
	}
	return nil
}
