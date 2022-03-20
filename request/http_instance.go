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

package request

import (
	"time"

	"github.com/vicanso/go-axios"
)

var insList = map[string]*axios.Instance{}

type InstanceStats struct {
	Name           string `json:"name"`
	MaxConcurrency int    `json:"maxConcurrency"`
	Concurrency    int    `json:"concurrency"`
}

// NewHTTP 新建实例
func NewHTTP(serviceName, baseURL string, timeout time.Duration) *axios.Instance {
	ins := axios.NewInstance(&axios.InstanceConfig{
		EnableTrace: true,
		Timeout:     timeout,
		OnError:     newOnError(serviceName),
		OnDone:      newOnDone(serviceName),
		BaseURL:     baseURL,
		ResponseInterceptors: []axios.ResponseInterceptor{
			newConvertResponseToError(),
		},
	})
	insList[serviceName] = ins
	return ins
}

// GetHTTPStats get http instance stats
func GetHTTPStats() []*InstanceStats {
	statsList := make([]*InstanceStats, len(insList))
	index := 0
	for name, ins := range insList {
		stats := InstanceStats{
			Name:           name,
			MaxConcurrency: int(ins.Config.MaxConcurrency),
			Concurrency:    int(ins.GetConcurrency()),
		}
		statsList[index] = &stats
		index++
	}
	return statsList
}

// UpdateConcurrencyLimit update the concurrency limit for instance
func UpdateConcurrencyLimit(limits map[string]int) {
	for name, ins := range insList {
		v := limits[name]
		limit := int32(v)
		if ins.Config.MaxConcurrency != limit {
			ins.SetMaxConcurrency(limit)
		}
	}
}
