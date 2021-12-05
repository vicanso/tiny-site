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
	"net"
	"net/http"

	"github.com/shirou/gopsutil/v3/process"
	performance "github.com/vicanso/go-performance"
)

var httpServerConnStats = performance.NewHttpServerConnStats()
var requestProcessConcurrency = performance.NewConcurrency()

// GetHTTPServerConnState get http server on conn state function
func GetHTTPServerConnState() func(net.Conn, http.ConnState) {
	return httpServerConnStats.ConnState
}

// dec http connection processing
func DecHTTPConnProcessing() {
	httpServerConnStats.DecProcessing()
}

type (
	// Performance 应用性能指标
	Performance struct {
		Concurrency           int32 `json:"concurrency"`
		RequestProcessedTotal int64 `json:"requestProcessedTotal"`
		performance.CPUMemory
		HTTPServerConnStats *performance.ConnStats        `json:"httpServerConnStats"`
		IOCountersStat      *process.IOCountersStat       `json:"ioCountersStat"`
		ConnStat            *performance.ConnectionsCount `json:"connStat"`
		NumCtxSwitchesStat  *process.NumCtxSwitchesStat   `json:"numCtxSwitchesStat"`
		PageFaultsStat      *process.PageFaultsStat       `json:"pageFaultsStat"`
		NumFds              int                           `json:"numFds"`
		OpenFilesStats      []process.OpenFilesStat       `json:"openFilesStats"`
	}
)

// GetPerformance 获取应用性能指标
func GetPerformance(ctx context.Context) *Performance {
	ioCountersStat, _ := performance.IOCounters(ctx)
	connStat, _ := performance.ConnectionsStat(ctx)
	numCtxSwitchesStat, _ := performance.NumCtxSwitches(ctx)
	numFds, _ := performance.NumFds(ctx)
	pageFaults, _ := performance.PageFaults(ctx)
	openFilesStats, _ := performance.OpenFiles(ctx)
	httpServerConnStats := httpServerConnStats.Stats()
	return &Performance{
		Concurrency:           GetConcurrency(),
		RequestProcessedTotal: requestProcessConcurrency.Total(),
		CPUMemory:             performance.CurrentCPUMemory(ctx),
		HTTPServerConnStats:   &httpServerConnStats,
		IOCountersStat:        ioCountersStat,
		ConnStat:              connStat,
		NumCtxSwitchesStat:    numCtxSwitchesStat,
		NumFds:                int(numFds),
		PageFaultsStat:        pageFaults,
		OpenFilesStats:        openFilesStats,
	}
}

// IncreaseConcurrency 当前并发请求+1
func IncreaseConcurrency() int32 {
	return requestProcessConcurrency.Inc()
}

// DecreaseConcurrency 当前并发请求-1
func DecreaseConcurrency() int32 {
	return requestProcessConcurrency.Dec()
}

// GetConcurrency 获取当前并发请求
func GetConcurrency() int32 {
	return requestProcessConcurrency.Current()
}

// UpdateCPUUsage 更新cpu使用率
func UpdateCPUUsage(ctx context.Context) error {
	return performance.UpdateCPUUsage(ctx)
}
