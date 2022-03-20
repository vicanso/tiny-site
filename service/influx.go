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

package service

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/vicanso/tiny-site/cache"
	"github.com/vicanso/tiny-site/helper"
)

type (
	InfluxSrv struct {
		db *helper.InfluxDB
	}
)

// 缓存的数据
type (
	// 缓存的values
	fluxCacheValues struct {
		Values []string `json:"values"`
	}
)

var defaultInfluxSrv = mustNewInfluxSrv()
var ignoreFields = "_start _stop _field _measurement"

// 缓存flux的tags tag-values等
var fluxCache = cache.NewMultilevelCache(10, 5*time.Minute, "flux:")

// mustNewInfluxSrv 创建新的influx服务
func mustNewInfluxSrv() *InfluxSrv {
	return &InfluxSrv{
		db: helper.GetInfluxDB(),
	}
}

// GetInfluxSrv 获取默认的influxdb服务
func GetInfluxSrv() *InfluxSrv {
	return defaultInfluxSrv
}

// ListTagValue list value of tag
func (srv *InfluxSrv) ListTagValue(ctx context.Context, measurement, tag string) (values []string, err error) {
	if srv.db == nil {
		return
	}
	// 优先取缓存
	key := fmt.Sprintf("tagValues:%s:%s", measurement, tag)
	result := fluxCacheValues{}
	// 忽略获取失败
	_ = fluxCache.Get(ctx, key, &result)
	if len(result.Values) != 0 {
		values = result.Values
		return
	}
	query := fmt.Sprintf(`import "influxdata/influxdb/schema"
schema.measurementTagValues(
	bucket: "%s",
	measurement: "%s",
	tag: "%s"
)`, srv.db.GetBucket(), measurement, tag)
	items, err := srv.db.Query(ctx, query)
	if err != nil {
		return
	}
	for _, item := range items {
		v, ok := item["_value"]
		if !ok {
			continue
		}
		value, ok := v.(string)
		if !ok {
			continue
		}
		values = append(values, value)
	}
	sort.Strings(values)
	if len(values) != 0 {
		result.Values = values
		_ = fluxCache.Set(ctx, key, &result)
	}
	return
}

// ListTag returns the tag list of measurement
func (srv *InfluxSrv) ListTag(ctx context.Context, measurement string) (tags []string, err error) {
	if srv.db == nil {
		return
	}
	// 优先取缓存
	key := fmt.Sprintf("tags:%s", measurement)
	result := fluxCacheValues{}
	_ = fluxCache.Get(ctx, key, &result)
	if len(result.Values) != 0 {
		tags = result.Values
		return
	}
	query := fmt.Sprintf(`import "influxdata/influxdb/schema"
schema.measurementTagKeys(
	bucket: "%s",
	measurement: "%s"
)`, srv.db.GetBucket(), measurement)
	items, err := srv.db.Query(ctx, query)
	if err != nil {
		return
	}
	for _, item := range items {
		v, ok := item["_value"]
		if !ok {
			continue
		}
		tag, ok := v.(string)
		if !ok {
			continue
		}
		if strings.Contains(ignoreFields, tag) {
			continue
		}
		tags = append(tags, tag)
	}
	if len(tags) != 0 {
		result.Values = tags
		_ = fluxCache.Set(ctx, key, &result)
	}
	return
}

func (srv *InfluxSrv) Query(ctx context.Context, query string) (items []map[string]interface{}, err error) {
	if srv.db == nil {
		return
	}
	query = fmt.Sprintf(`from(bucket: "%s")
`, srv.db.GetBucket()) + query
	return srv.db.Query(ctx, query)
}

// ListField return the fields of measurement
func (srv *InfluxSrv) ListField(ctx context.Context, measurement, duration string) (fields []string, err error) {
	if srv.db == nil {
		return
	}
	// 优先取缓存
	key := fmt.Sprintf("fields:%s:%s", measurement, duration)
	result := fluxCacheValues{}
	_ = fluxCache.Get(ctx, key, &result)
	if len(result.Values) != 0 {
		fields = result.Values
		return
	}

	query := fmt.Sprintf(`import "influxdata/influxdb/schema"
schema.measurementFieldKeys(
	bucket: "%s",
	measurement: "%s",
	start: %s
)`, srv.db.GetBucket(), measurement, duration)
	items, err := srv.db.Query(ctx, query)
	if err != nil {
		return
	}
	for _, item := range items {
		v, ok := item["_value"]
		if !ok {
			continue
		}
		field, ok := v.(string)
		if !ok {
			continue
		}
		if strings.Contains(ignoreFields, field) {
			continue
		}
		fields = append(fields, field)
	}
	if len(fields) != 0 {
		result.Values = fields
		_ = fluxCache.Set(ctx, key, &result)
	}
	return
}

// Write 写入数据
func (srv *InfluxSrv) Write(measurement string, tags map[string]string, fields map[string]interface{}, ts ...time.Time) {
	if srv.db == nil {
		return
	}
	srv.db.Write(measurement, tags, fields, ts...)
}
