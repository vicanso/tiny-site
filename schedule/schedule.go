// Copyright 2019 tree xie
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
	"time"

	"github.com/vicanso/tiny-site/log"
	"github.com/vicanso/tiny-site/service"

	"go.uber.org/zap"
)

func init() {
	go initRedisCheckTicker()
	go initConfigurationRefreshTicker()
	// go initInfluxdbCheckTicker()
	// go initRouterConfigRefreshTicker()
}

func runTicker(ticker *time.Ticker, message string, do func() error, restart func()) {
	defer func() {
		if r := recover(); r != nil {
			err, _ := r.(error)
			log.Default().DPanic(message+" panic",
				zap.Error(err),
			)
		}
		// 如果退出了，重新启动
		go restart()
	}()
	for range ticker.C {
		err := do()
		// TODO 检测不通过时，发送告警
		if err != nil {
			log.Default().Error(message+" fail",
				zap.Error(err),
			)
		}
	}
}

func initRedisCheckTicker() {
	// 每一分钟检测一次
	ticker := time.NewTicker(60 * time.Second)
	runTicker(ticker, "redis check", func() error {
		err := service.RedisPing()
		return err
	}, initRedisCheckTicker)
}

func initConfigurationRefreshTicker() {
	// 每一分钟更新一次
	configSrv := new(service.ConfigurationSrv)
	ticker := time.NewTicker(60 * time.Second)
	runTicker(ticker, "configuration refresh", func() error {
		err := configSrv.Refresh()
		return err
	}, initConfigurationRefreshTicker)
}
