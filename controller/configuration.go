// Copyright 2020 tree xie
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

// 应用相关配置，包括IP拦截、路由mock、路由并发限制等配置信息

package controller

import (
	"context"
	"time"

	"github.com/vicanso/elton"
	"github.com/vicanso/hes"
	"github.com/vicanso/tiny-site/cs"
	"github.com/vicanso/tiny-site/ent"
	"github.com/vicanso/tiny-site/ent/configuration"
	confSchema "github.com/vicanso/tiny-site/ent/configuration"
	"github.com/vicanso/tiny-site/router"
	"github.com/vicanso/tiny-site/schema"
	"github.com/vicanso/tiny-site/service"
)

type configurationCtrl struct{}

// 响应相关定义
type (
	// configurationListResp 配置列表响应
	configurationListResp struct {
		Configurations []*ent.Configuration `json:"configurations"`
		Count          int                  `json:"count"`
	}
)

// 参数相关定义
type (

	// configurationAddParams 添加配置参数
	configurationAddParams struct {
		Name        string              `json:"name" validate:"required,xConfigurationName"`
		Category    confSchema.Category `json:"category" validate:"required,xConfigurationCategory"`
		Status      schema.Status       `json:"status" validate:"required,xStatus"`
		Data        string              `json:"data" validate:"required,xConfigurationData"`
		StartedAt   time.Time           `json:"startedAt"`
		EndedAt     time.Time           `json:"endedAt"`
		Description string              `json:"description"`
	}
	// configurationUpdateParams 更新配置参数
	configurationUpdateParams struct {
		Name        string              `json:"name" validate:"omitempty,xConfigurationName"`
		Status      schema.Status       `json:"status" validate:"omitempty,xStatus"`
		Category    confSchema.Category `json:"category" validate:"omitempty,xConfigurationCategory"`
		Data        string              `json:"data" validate:"omitempty,xConfigurationData"`
		StartedAt   time.Time           `json:"startedAt"`
		EndedAt     time.Time           `json:"endedAt"`
		Description string              `json:"description"`
	}

	// configurationListParmas 配置查询参数
	configurationListParmas struct {
		listParams

		Name     string              `json:"name" validate:"omitempty,xConfigurationName"`
		Category confSchema.Category `json:"category" validate:"omitempty,xConfigurationCategory"`
	}
)

const (
	errConfigurationCategory = "configuration"
)

func init() {
	g := router.NewGroup(
		"/configurations",
		loadUserSession,
		shouldBeSu,
	)
	ctrl := configurationCtrl{}

	// 查询配置
	g.GET(
		"/v1",
		ctrl.list,
	)

	// 添加配置
	g.POST(
		"/v1",
		newTrackerMiddleware(cs.ActionConfigurationAdd),
		ctrl.add,
	)

	// 获取当前有效配置
	g.GET(
		"/v1/current-valid",
		ctrl.getCurrentValid,
	)

	// 更新配置
	g.PATCH(
		"/v1/{id}",
		newTrackerMiddleware(cs.ActionConfigurationUpdate),
		ctrl.update,
	)

	// 查询单个配置
	g.GET(
		"/v1/{id}",
		ctrl.findByID,
	)
}

// validateBeforeSave 保存前校验
func (params *configurationAddParams) validateBeforeSave(ctx context.Context) error {
	// schema中有唯一限制，也可不校验
	exists, err := getConfigurationClient().Query().
		Where(configuration.Name(params.Name)).
		Exist(ctx)
	if err != nil {
		return err
	}
	if exists {
		return hes.New("该配置已存在", errConfigurationCategory)
	}
	return nil
}

// save 保存配置
func (params *configurationAddParams) save(ctx context.Context, owner string) (*ent.Configuration, error) {
	err := params.validateBeforeSave(ctx)
	if err != nil {
		return nil, err
	}
	return getConfigurationClient().Create().
		SetName(params.Name).
		SetStatus(params.Status).
		SetCategory(params.Category).
		SetData(params.Data).
		SetOwner(owner).
		SetStartedAt(params.StartedAt).
		SetEndedAt(params.EndedAt).
		SetDescription(params.Description).
		Save(ctx)
}

// where 将查询条件中的参数转换为对应的where条件
func (params *configurationListParmas) where(query *ent.ConfigurationQuery) *ent.ConfigurationQuery {
	if params.Name != "" {
		query.Where(configuration.Name(params.Name))
	}
	if params.Category != "" {
		query.Where(configuration.CategoryEQ(params.Category))
	}
	return query
}

// queryAll 查询配置列表
func (params *configurationListParmas) queryAll(ctx context.Context) ([]*ent.Configuration, error) {
	query := getConfigurationClient().Query()

	query.Limit(params.GetLimit()).
		Offset(params.GetOffset()).
		Order(params.GetOrders()...)
	params.where(query)

	return query.All(ctx)
}

// count 计算总数
func (params *configurationListParmas) count(ctx context.Context) (int, error) {
	query := getConfigurationClient().Query()

	params.where(query)

	return query.Count(ctx)
}

// update 更新配置信息
func (params *configurationUpdateParams) updateOneID(ctx context.Context, id int) (*ent.Configuration, error) {
	updateOne := getConfigurationClient().
		UpdateOneID(id)
	if !params.StartedAt.IsZero() {
		updateOne = updateOne.SetStartedAt(params.StartedAt)
	}
	if !params.EndedAt.IsZero() {
		updateOne = updateOne.SetEndedAt(params.EndedAt)
	}
	if params.Name != "" {
		updateOne = updateOne.SetName(params.Name)
	}

	if params.Status != 0 {
		updateOne = updateOne.SetStatus(params.Status)
	}
	if params.Category != "" {
		updateOne = updateOne.SetCategory(params.Category)
	}
	if params.Data != "" {
		updateOne = updateOne.SetData(params.Data)
	}
	if params.Description != "" {
		updateOne = updateOne.SetDescription(params.Description)
	}
	return updateOne.Save(ctx)
}

// add 添加配置
func (*configurationCtrl) add(c *elton.Context) error {
	params := configurationAddParams{}
	err := validateBody(c, &params)
	if err != nil {
		return err
	}
	us := getUserSession(c)
	configuration, err := params.save(c.Context(), us.MustGetInfo().Account)
	if err != nil {
		return err
	}
	c.Created(configuration)
	return nil
}

// list 查询配置列表
func (*configurationCtrl) list(c *elton.Context) error {
	params := configurationListParmas{}
	err := validateQuery(c, &params)
	if err != nil {
		return err
	}
	count := -1
	if params.ShouldCount() {
		count, err = params.count(c.Context())
		if err != nil {
			return err
		}
	}
	configurations, err := params.queryAll(c.Context())
	if err != nil {
		return err
	}
	c.Body = &configurationListResp{
		Count:          count,
		Configurations: configurations,
	}
	return nil
}

// update 更新配置信息
func (*configurationCtrl) update(c *elton.Context) error {
	id, err := getIDFromParams(c)
	if err != nil {
		return err
	}
	params := configurationUpdateParams{}
	err = validateBody(c, &params)
	if err != nil {
		return err
	}
	configuration, err := params.updateOneID(c.Context(), id)
	if err != nil {
		return err
	}

	c.Body = configuration
	return nil
}

// findByID 通过id查询
func (*configurationCtrl) findByID(c *elton.Context) error {
	id, err := getIDFromParams(c)
	if err != nil {
		return err
	}
	configuration, err := getConfigurationClient().Get(c.Context(), id)
	if err != nil {
		return err
	}
	c.Body = configuration
	return nil
}

// getCurrentValid 获取当前有效配置
func (*configurationCtrl) getCurrentValid(c *elton.Context) error {
	c.Body = service.GetCurrentValidConfiguration()
	return nil
}
