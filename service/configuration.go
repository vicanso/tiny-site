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

package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/vicanso/elton"
	"github.com/vicanso/tiny-site/config"
	"github.com/vicanso/tiny-site/email"
	"github.com/vicanso/tiny-site/ent"
	"github.com/vicanso/tiny-site/ent/configuration"
	"github.com/vicanso/tiny-site/helper"
	"github.com/vicanso/tiny-site/interceptor"
	"github.com/vicanso/tiny-site/log"
	"github.com/vicanso/tiny-site/request"
	routerconcurrency "github.com/vicanso/tiny-site/router_concurrency"
	routermock "github.com/vicanso/tiny-site/router_mock"
	"github.com/vicanso/tiny-site/schema"
	"github.com/vicanso/tiny-site/util"
	"go.uber.org/atomic"
)

// ConfigurationSrv 配置的相关函数
type ConfigurationSrv struct{}

// 配置数据
type (

	// CurrentValidConfiguration 当前有效配置
	CurrentValidConfiguration struct {
		UpdatedAt         time.Time                        `json:"updatedAt"`
		MockTime          string                           `json:"mockTime"`
		IPBlockList       []string                         `json:"ipBlockList"`
		SignedKeys        []string                         `json:"signedKeys"`
		RouterConcurrency map[string]uint32                `json:"routerConcurrency"`
		RouterMock        map[string]routermock.RouterMock `json:"routerMock"`
		Limits            map[string]int                   `json:"limits"`
	}
	// RequestLimitConfiguration HTTP请求实例并发限制
	RequestLimitConfiguration struct {
		Name string `json:"name"`
		Max  int    `json:"max"`
	}
)

var (
	sessionSignedKeys = new(elton.RWMutexSignedKeys)

	// 当前请求实例限制
	currentLimits = atomic.Value{}
)

// 配置刷新时间
var configurationRefreshedAt time.Time
var sessionConfig = config.MustGetSessionConfig()

func init() {
	// session中用于cookie的signed keys
	sessionSignedKeys.SetKeys(sessionConfig.Keys)
}

// GetSignedKeys 获取用于cookie加密的key列表
func GetSignedKeys() elton.SignedKeysGenerator {
	return sessionSignedKeys
}

// GetCurrentValidConfiguration 获取当前有效配置
func GetCurrentValidConfiguration() *CurrentValidConfiguration {
	result := &CurrentValidConfiguration{
		UpdatedAt:         configurationRefreshedAt,
		MockTime:          util.GetMockTime(),
		IPBlockList:       GetIPBlockList(),
		SignedKeys:        sessionSignedKeys.GetKeys(),
		RouterConcurrency: routerconcurrency.List(),
		RouterMock:        routermock.List(),
	}
	value := currentLimits.Load()
	if value != nil {
		limits, _ := value.(map[string]int)
		result.Limits = limits
	}
	return result
}

// available 获取可用的配置
func (*ConfigurationSrv) available() ([]*ent.Configuration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	now := time.Now()
	return helper.EntGetClient().Configuration.Query().
		Where(configuration.Status(schema.StatusEnabled)).
		Where(configuration.StartedAtLT(now)).
		Where(configuration.EndedAtGT(now)).
		Order(ent.Desc(configuration.FieldUpdatedAt)).
		All(ctx)
}

// Refresh 刷新配置
func (srv *ConfigurationSrv) Refresh() error {
	configs, err := srv.available()
	if err != nil {
		return err
	}
	configurationRefreshedAt = time.Now()
	var mockTimeConfig *ent.Configuration
	routerConcurrencyConfigs := make([]string, 0)
	routerConfigs := make([]string, 0)
	var signedKeys []string
	blockIPList := make([]string, 0)

	mailList := make(map[string]string)

	httpInterceptors := make([]string, 0)

	requestLimitConfigs := make(map[string]int)
	for _, item := range configs {
		switch item.Category {
		case schema.ConfigurationCategoryMockTime:
			// 由于排序是按更新时间，因此取最新的记录
			if mockTimeConfig == nil {
				mockTimeConfig = item
			}
		case schema.ConfigurationCategoryBlockIP:
			blockIPList = append(blockIPList, item.Data)
		case schema.ConfigurationCategorySignedKey:
			// 按更新时间排序，因此如果已获取则不需要再更新
			if len(signedKeys) == 0 {
				signedKeys = strings.Split(item.Data, ",")
			}
		case schema.ConfigurationCategoryRouterConcurrency:
			routerConcurrencyConfigs = append(routerConcurrencyConfigs, item.Data)
		case schema.ConfigurationCategoryRouter:
			routerConfigs = append(routerConfigs, item.Data)
		case schema.ConfigurationCategoryRequestConcurrency:
			c := RequestLimitConfiguration{}
			err := json.Unmarshal([]byte(item.Data), &c)
			if err != nil {
				log.Error(context.Background()).
					Err(err).
					Msg("request limit config is invalid")
				email.AlarmError(context.Background(), "request limit config is invalid:"+err.Error())
			}
			if c.Name != "" {
				requestLimitConfigs[c.Name] = c.Max
			}
		case schema.ConfigurationCategoryEmail:
			mailList[item.Name] = item.Data
		case schema.ConfigurationHTTPServerInterceptor:
			httpInterceptors = append(httpInterceptors, item.Data)
		}
	}

	// 如果未配置mock time，则设置为空
	if mockTimeConfig == nil {
		util.SetMockTime("")
	} else {
		util.SetMockTime(mockTimeConfig.Data)
	}

	// 如果数据库中未配置，则使用默认配置
	if len(signedKeys) == 0 {
		sessionSignedKeys.SetKeys(sessionConfig.Keys)
	} else {
		sessionSignedKeys.SetKeys(signedKeys)
	}

	// 更新router configs
	routermock.Update(routerConfigs)

	// 重置IP拦截列表
	err = ResetIPBlocker(blockIPList)
	if err != nil {
		log.Error(context.Background()).
			Err(err).
			Msg("reset ip blocker fail")
	}

	// 重置路由并发限制
	routerconcurrency.Update(routerConcurrencyConfigs)

	// 更新HTTP请求实例并发限制
	currentLimits.Store(requestLimitConfigs)
	request.UpdateConcurrencyLimit(requestLimitConfigs)

	email.Update(mailList)

	// 更新拦截配置
	interceptor.UpdateHTTPInterceptors(httpInterceptors)

	return nil
}
