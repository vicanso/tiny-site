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

package schedule

import (
	"context"
	"strconv"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/vicanso/tiny-site/cs"
	"github.com/vicanso/tiny-site/email"
	"github.com/vicanso/tiny-site/helper"
	"github.com/vicanso/tiny-site/log"
	"github.com/vicanso/tiny-site/request"
	routerconcurrency "github.com/vicanso/tiny-site/router_concurrency"
	"github.com/vicanso/tiny-site/service"
	"github.com/vicanso/tiny-site/util"
)

type (
	taskFn      func() error
	statsTaskFn func() map[string]interface{}
)

const logCategory = "schedule"

func init() {
	c := cron.New()
	_, _ = c.AddFunc("@every 1m", redisPing)
	_, _ = c.AddFunc("@every 1m", entPing)
	_, _ = c.AddFunc("@every 1m", influxdbPing)
	_, _ = c.AddFunc("@every 1m", configRefresh)
	_, _ = c.AddFunc("@every 1m", redisStats)
	_, _ = c.AddFunc("@every 1m", entStats)
	_, _ = c.AddFunc("@every 30s", cpuUsageStats)
	_, _ = c.AddFunc("@every 1m", performanceStats)
	_, _ = c.AddFunc("@every 1m", httpInstanceStats)
	_, _ = c.AddFunc("@every 1m", routerConcurrencyStats)
	// 如果是开发环境，则不执行定时任务
	if util.IsDevelopment() {
		return
	}
	c.Start()
}

func doTask(desc string, fn taskFn) {
	startedAt := time.Now()
	err := fn()
	use := time.Since(startedAt)
	if err != nil {
		log.Error(context.Background()).
			Str("category", logCategory).
			Dur("use", use).
			Err(err).
			Msg(desc + " fail")
		email.AlarmError(context.Background(), desc+" fail, "+err.Error())
	} else {
		log.Info(context.Background()).
			Str("category", logCategory).
			Dur("use", use).
			Msg(desc + " success")
	}
}

func doStatsTask(desc string, fn statsTaskFn) {
	startedAt := time.Now()
	stats := fn()
	log.Info(context.Background()).
		Str("category", logCategory).
		Dur("use", time.Since(startedAt)).
		Dict("stats", zerolog.Dict().Fields(stats)).
		Msg("")
}

func redisPing() {
	doTask("redis ping", helper.RedisPing)
}

func configRefresh() {
	configSrv := new(service.ConfigurationSrv)
	doTask("config refresh", configSrv.Refresh)
}

func redisStats() {
	doStatsTask("redis stats", func() map[string]interface{} {
		// 统计中除了redis数据库的统计，还有当前实例的统计指标，因此所有实例都会写入统计
		stats := helper.RedisStats()
		helper.GetInfluxDB().Write(cs.MeasurementRedisStats, nil, stats)
		return stats
	})
}

func entPing() {
	doTask("ent ping", helper.EntPing)
}

// entStats ent的性能统计
func entStats() {
	doStatsTask("ent stats", func() map[string]interface{} {
		stats := helper.EntGetStats()
		helper.GetInfluxDB().Write(cs.MeasurementEntStats, nil, stats)
		return stats
	})
}

// cpuUsageStats cpu使用率
func cpuUsageStats() {
	doTask("update cpu usage", func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		return service.UpdateCPUUsage(ctx)
	})
}

// prevMemFrees 上一次 memory objects 释放的数量
var prevMemFrees uint64

// 上一次memory objects的申请数量
var prevMemMallocs uint64

// prevNumGC 上一次 gc 的次数
var prevNumGC uint32

// prevPauseTotal 上一次 pause 的总时长
var prevPauseTotal time.Duration

// 上一次的 io counter记录
var prevIOCountersStat = &process.IOCountersStat{}

// 上一次context切换记录
var prevNumCtxSwitchesStat = &process.NumCtxSwitchesStat{}

// 上一次的page faults记录
var prevPageFaultsStat = &process.PageFaultsStat{}

const mb = 1024 * 1024

// performanceStats 系统性能
func performanceStats() {
	doStatsTask("performance stats", func() map[string]interface{} {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		data := service.GetPerformance(ctx)
		fields := map[string]interface{}{
			cs.FieldGoMaxProcs:   data.GoMaxProcs,
			cs.FieldThreadCount:  int(data.ThreadCount),
			cs.FieldRoutineCount: data.RoutineCount,
			cs.FieldProcessing:   int(data.Concurrency),

			// CPU使用率相关
			cs.FieldCPUUsage:     int(data.CPUUsage),
			cs.FieldCPUUser:      data.CPUUser,
			cs.FieldCPUIdle:      data.CPUIdle,
			cs.FieldCPUNice:      data.CPUNice,
			cs.FieldCPUIowait:    data.CPUIowait,
			cs.FieldCPUIrq:       data.CPUIrq,
			cs.FieldCPUSoftirq:   data.CPUSoftirq,
			cs.FieldCPUSteal:     data.CPUSteal,
			cs.FieldCPUGuest:     data.CPUGuest,
			cs.FieldCPUGuestNice: data.CPUGuestNice,

			// 内存使用相关
			cs.FieldMemAlloc:        data.MemAlloc,
			cs.FieldMemTotalAlloc:   data.MemTotalAlloc,
			cs.FieldMemSys:          data.MemSys,
			cs.FieldMemLookups:      data.MemLookups,
			cs.FieldMemMallocs:      int(data.MemMallocs - prevMemMallocs),
			cs.FieldMemFrees:        int(data.MemFrees - prevMemFrees),
			cs.FieldMemHeapAlloc:    data.MemHeapAlloc,
			cs.FieldMemHeapSys:      data.MemHeapSys,
			cs.FieldMemHeapIdle:     data.MemHeapIdle,
			cs.FieldMemHeapInuse:    data.MemHeapInuse,
			cs.FieldMemHeapReleased: data.MemHeapReleased,
			cs.FieldMemHeapObjects:  data.MemHeapObjects,
			cs.FieldMemStackInuse:   data.MemStackInuse,
			cs.FieldMemStackSys:     data.MemStackSys,
			cs.FieldMemMSpanInuse:   data.MemMSpanInuse,
			cs.FieldMemMSpanSys:     data.MemMSpanSys,
			cs.FieldMemMCacheInuse:  data.MemMCacheInuse,
			cs.FieldMemMCacheSys:    data.MemMCacheSys,
			cs.FieldMemBuckHashSys:  data.MemBuckHashSys,

			cs.FieldMemGCSys:    data.MemGCSys,
			cs.FieldMemOtherSys: data.MemOtherSys,

			cs.FieldNumGC:   int(data.NumGC - prevNumGC),
			cs.FieldPauseMS: int((data.PauseTotalNs - prevPauseTotal).Milliseconds()),

			cs.FieldConnProcessing:   int(data.HTTPServerConnStats.ConnProcessing),
			cs.FieldConnAlive:        int(data.HTTPServerConnStats.ConnAlive),
			cs.FieldConnCreatedCount: int(data.HTTPServerConnStats.ConnCreatedCount),
		}
		// io 相关统计
		if data.IOCountersStat != nil {
			readCount := data.IOCountersStat.ReadCount - prevIOCountersStat.ReadCount
			writeCount := data.IOCountersStat.WriteCount - prevIOCountersStat.WriteCount
			readBytes := data.IOCountersStat.ReadBytes - prevIOCountersStat.ReadBytes
			writeBytes := data.IOCountersStat.WriteBytes - prevIOCountersStat.WriteBytes
			prevIOCountersStat = data.IOCountersStat
			fields[cs.FieldIOReadCount] = readCount
			fields[cs.FieldIOWriteCount] = writeCount
			fields[cs.FieldIOReadMBytes] = readBytes / mb
			fields[cs.FieldIOWriteMBytes] = writeBytes / mb
		}

		// 网络相关
		if data.ConnStat != nil {
			count := make(map[string]string)

			for k, v := range data.ConnStat.Status {
				count[k] = strconv.Itoa(v)
			}
			for k, v := range data.ConnStat.RemoteAddr {
				count[k] = strconv.Itoa(v)
			}

			log.Info(context.Background()).
				Str("category", "connStat").
				Dict("count", log.Struct(count)).
				Msg("")
			fields[cs.FieldConnTotal] = data.ConnStat.Count
		}
		// context切换相关
		if data.NumCtxSwitchesStat != nil {
			fields[cs.FieldCtxSwitchesVoluntary] = data.NumCtxSwitchesStat.Voluntary - prevNumCtxSwitchesStat.Voluntary
			fields[cs.FieldCtxSwitchesInvoluntary] = data.NumCtxSwitchesStat.Involuntary - prevNumCtxSwitchesStat.Involuntary
			prevNumCtxSwitchesStat = data.NumCtxSwitchesStat
		}

		// page fault相关
		if data.PageFaultsStat != nil {
			fields[cs.FieldPageFaultMinor] = data.PageFaultsStat.MinorFaults - prevPageFaultsStat.MinorFaults
			fields[cs.FieldPageFaultMajor] = data.PageFaultsStat.MajorFaults - prevPageFaultsStat.MajorFaults
			fields[cs.FieldPageFaultChildMinor] = data.PageFaultsStat.ChildMinorFaults - prevPageFaultsStat.ChildMinorFaults
			fields[cs.FieldPageFaultChildMajor] = data.PageFaultsStat.ChildMajorFaults - prevPageFaultsStat.ChildMajorFaults

			prevPageFaultsStat = data.PageFaultsStat
		}

		// fd 相关
		fields[cs.FieldNumFds] = data.NumFds

		// open files的统计
		if len(data.OpenFilesStats) != 0 {
			log.Info(context.Background()).
				Str("category", "openFiles").
				Dict("stat", log.Struct(data.OpenFilesStats)).
				Msg("")
		}

		prevMemMallocs = data.MemMallocs
		prevMemFrees = data.MemFrees
		prevNumGC = data.NumGC
		prevPauseTotal = data.PauseTotalNs

		helper.GetInfluxDB().Write(cs.MeasurementPerformance, nil, fields)
		return fields
	})
}

// httpInstanceStats http instance stats
func httpInstanceStats() {
	doStatsTask("http instance stats", func() map[string]interface{} {
		fields := make(map[string]interface{})
		statsList := request.GetHTTPStats()
		for _, stats := range statsList {
			helper.GetInfluxDB().Write(cs.MeasurementHTTPInstanceStats, map[string]string{
				cs.TagService: stats.Name,
			}, map[string]interface{}{
				cs.FieldMaxConcurrency: stats.MaxConcurrency,
				cs.FieldProcessing:     stats.Concurrency,
			})
			fields[stats.Name+":"+cs.FieldMaxConcurrency] = stats.MaxConcurrency
			fields[stats.Name+":"+cs.FieldProcessing] = stats.Concurrency
		}
		return fields
	})
}

// influxdbPing influxdb ping
func influxdbPing() {
	doTask("influxdb ping", helper.GetInfluxDB().Health)
}

// routerConcurrencyStats router concurrency stats
func routerConcurrencyStats() {
	doStatsTask("router concurrency stats", func() map[string]interface{} {
		result := routerconcurrency.GetLimiter().GetStats()
		fields := make(map[string]interface{})

		influxSrv := helper.GetInfluxDB()
		for key, value := range result {
			// 如果并发为0，则不记录
			if value == 0 {
				continue
			}
			fields[key] = value
			influxSrv.Write(cs.MeasurementRouterConcurrency, map[string]string{
				cs.TagRoute: key,
			}, map[string]interface{}{
				cs.FieldCount: int(value),
			})
		}
		return fields
	})
}
