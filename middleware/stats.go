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

package middleware

import (
	"github.com/dustin/go-humanize"
	"github.com/vicanso/elton"
	M "github.com/vicanso/elton/middleware"
	"github.com/vicanso/tiny-site/cs"
	"github.com/vicanso/tiny-site/helper"
	"github.com/vicanso/tiny-site/log"
	"github.com/vicanso/tiny-site/util"
)

func NewStats() elton.Handler {
	return M.NewStats(M.StatsConfig{
		OnStats: func(info *M.StatsInfo, c *elton.Context) {
			// ping 的日志忽略
			if info.URI == "/ping" {
				return
			}
			// TODO 如果需要可以从session中获取账号
			// 此处记录的session id，需要从客户登录记录中获取对应的session id
			// us := service.NewUserSession(c)
			sid := util.GetSessionID(c)
			// 由客户端设置的uuid
			// zap.String("uuid", c.GetRequestHeader("X-UUID")),
			log.Info(c.Context()).
				Str("category", "accessLog").
				Str("ip", info.IP).
				Str("sid", sid).
				Str("method", info.Method).
				Str("route", info.Route).
				Str("uri", info.URI).
				Int("status", info.Status).
				Uint32("connecting", info.Connecting).
				Str("latency", info.Latency.String()).
				Str("size", humanize.Bytes(uint64(info.Size))).
				Int("bytes", info.Size).
				Msg("")

			tags := map[string]string{
				cs.TagMethod: info.Method,
				cs.TagRoute:  info.Route,
			}
			fields := map[string]interface{}{
				cs.FieldIP:         info.IP,
				cs.FieldSID:        sid,
				cs.FieldURI:        info.URI,
				cs.FieldStatus:     info.Status,
				cs.FieldLatency:    int(info.Latency.Milliseconds()),
				cs.FieldSize:       info.Size,
				cs.FieldProcessing: info.Connecting,
			}
			helper.GetInfluxDB().Write(cs.MeasurementHTTPStats, tags, fields)
		},
	})
}
