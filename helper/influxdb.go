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

package helper

import (
	"context"
	"errors"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	influxdbAPI "github.com/influxdata/influxdb-client-go/v2/api"
	influxdbDomain "github.com/influxdata/influxdb-client-go/v2/domain"
	"github.com/vicanso/tiny-site/config"
	"github.com/vicanso/tiny-site/cs"
	"github.com/vicanso/tiny-site/log"
	"go.uber.org/atomic"
)

type (
	InfluxDB struct {
		client            influxdb2.Client
		writer            influxdbAPI.WriteAPI
		config            *config.InfluxdbConfig
		writeCount        atomic.Int64
		writtingCount     atomic.Int32
		maxWrittingPoints int32
	}
)

var hostname, _ = os.Hostname()
var defaultInfluxDB = mustNewInfluxDB()

// mustNewInfluxDB 创建新的influx服务
func mustNewInfluxDB() *InfluxDB {
	influxdbConfig := config.MustGetInfluxdbConfig()
	if influxdbConfig.Disabled {

		return new(InfluxDB)
	}
	opts := influxdb2.DefaultOptions()
	// 设置批量提交的大小
	opts.SetBatchSize(influxdbConfig.BatchSize)
	// 如果定时提交间隔大于1秒，则设定定时提交间隔
	if influxdbConfig.FlushInterval > time.Second {
		v := influxdbConfig.FlushInterval / time.Millisecond
		opts.SetFlushInterval(uint(v))
	}
	opts.SetPrecision(time.Nanosecond)
	opts.SetUseGZip(influxdbConfig.Gzip)
	log.Info(context.Background()).
		Str("uri", influxdbConfig.URI).
		Str("org", influxdbConfig.Org).
		Str("bucket", influxdbConfig.Bucket).
		Uint("batchSize", influxdbConfig.BatchSize).
		Str("token", influxdbConfig.Token[:5]+"...").
		Dur("interval", influxdbConfig.FlushInterval).
		Msg("")

	c := influxdb2.NewClientWithOptions(influxdbConfig.URI, influxdbConfig.Token, opts)
	writer := c.WriteAPI(influxdbConfig.Org, influxdbConfig.Bucket)
	db := &InfluxDB{
		client:            c,
		writer:            writer,
		config:            influxdbConfig,
		maxWrittingPoints: int32(influxdbConfig.MaxWrittingPoints),
	}
	go newInfluxdbErrorLogger(writer, db)

	return db
}

// newInfluxdbErrorLogger 创建读取出错日志处理，需要注意此功能需要启用新的goroutine
func newInfluxdbErrorLogger(writer influxdbAPI.WriteAPI, db *InfluxDB) {
	for err := range writer.Errors() {
		log.Error(context.Background()).
			Str("category", "influxdbError").
			Err(err).
			Msg("")
		db.Write(cs.MeasurementException, map[string]string{
			cs.TagCategory: "influxdbError",
		}, map[string]interface{}{
			cs.FieldError: err.Error(),
		})
	}
}

// GetInfluxDB 获取默认的influxdb服务
func GetInfluxDB() *InfluxDB {
	return defaultInfluxDB
}

// Health check influxdb health
func (db *InfluxDB) Health() error {
	if db.client == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := db.client.Health(ctx)
	if err != nil {
		return err
	}
	if result.Status != influxdbDomain.HealthCheckStatusPass {
		return errors.New(string(result.Status))
	}
	return nil
}

func (db *InfluxDB) Query(ctx context.Context, query string) ([]map[string]interface{}, error) {
	if db.client == nil {
		return nil, nil
	}
	result, err := db.client.QueryAPI(db.config.Org).Query(ctx, query)
	if err != nil {
		return nil, err
	}
	items := make([]map[string]interface{}, 0)
	for result.Next() {
		items = append(items, result.Record().Values())
	}
	err = result.Err()
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *InfluxDB) GetAndResetWriteCount() int64 {
	return db.writeCount.Swap(0)
}

func (db *InfluxDB) GetWrittingCount() int32 {
	return db.writtingCount.Load()
}

// Write 写入数据
func (db *InfluxDB) Write(measurement string, tags map[string]string, fields map[string]interface{}, ts ...time.Time) {
	if db.writer == nil {
		return
	}
	db.writeCount.Inc()
	var now time.Time
	if len(ts) != 0 {
		now = ts[0]
	} else {
		now = time.Now()
	}
	if fields == nil {
		fields = make(map[string]interface{})
	}
	if hostname != "" && fields["hostname"] == nil {
		fields["hostname"] = hostname
	}
	value := db.writtingCount.Inc()
	defer db.writtingCount.Dec()
	// 由于write point有可能由于上一次batch提交处理中
	// 而新的batch又满了会导致卡住
	// 因此增加处理请求量超过20时，直接不再写统计
	if value > db.maxWrittingPoints {
		log.Error(context.Background()).
			Int32("count", value).
			Msg("too many points are waiting")
		return
	}
	db.writer.WritePoint(influxdb2.NewPoint(measurement, tags, fields, now))
}

// Close 关闭当前client
func (db *InfluxDB) Close() {
	if db.client == nil {
		return
	}
	db.client.Close()
}

// GetBucket 获取bucket
func (db *InfluxDB) GetBucket() string {
	return db.config.Bucket
}
