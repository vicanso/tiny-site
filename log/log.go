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

// 可通过zap.RegisterSink添加更多的sink实现不同方式的日志传输

package log

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/vicanso/tiny-site/cs"
	"github.com/vicanso/tiny-site/util"
	mask "github.com/vicanso/go-mask"
)

var enabledDebugLog = false
var defaultLogger = newLogger()

// 日志中值的最大长度
var logFieldValueMaxSize = 30
var logMask = mask.New(
	mask.RegExpOption(cs.MaskRegExp),
	mask.MaxLengthOption(logFieldValueMaxSize),
	mask.NotMaskRegExpOption(regexp.MustCompile(`stack`)),
)

type httpServerLogger struct{}

func (hsl *httpServerLogger) Write(p []byte) (int, error) {
	Info(context.Background()).
		Str("category", "httpServerLogger").
		Msg(string(p))
	return len(p), nil
}

type redisLogger struct{}

func (rl *redisLogger) Printf(ctx context.Context, format string, v ...interface{}) {
	Info(context.Background()).
		Str("category", "redisLogger").
		Msg(fmt.Sprintf(format, v...))
}

type entLogger struct{}

func (el *entLogger) Log(args ...interface{}) {
	Info(context.Background()).
		Msg(fmt.Sprint(args...))
}

// DebugEnabled 是否启用了debug日志
func DebugEnabled() bool {
	return enabledDebugLog
}

// newLogger 初始化logger
func newLogger() *zerolog.Logger {
	// 全局禁用sampling
	zerolog.DisableSampling(true)
	// 如果要节约日志空间，可以配置
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.999Z07:00"

	var l zerolog.Logger
	if util.IsDevelopment() {
		l = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
			With().
			Timestamp().
			Logger()
	} else {
		l = zerolog.New(os.Stdout).
			Level(zerolog.InfoLevel).
			With().
			Timestamp().
			Logger()
	}

	// 如果有配置指定日志级别，则以配置指定的输出
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		lv, _ := strconv.Atoi(logLevel)
		l = l.Level(zerolog.Level(lv))
		if logLevel != "" && lv <= 0 {
			enabledDebugLog = true
		}
	}

	return &l
}

func fillTraceInfos(ctx context.Context, e *zerolog.Event) *zerolog.Event {
	if ctx == nil {
		ctx = context.Background()
	}
	deviceID := util.GetDeviceID(ctx)
	if deviceID == "" {
		return e
	}
	return e.Str("deviceID", deviceID).
		Str("traceID", util.GetTraceID(ctx)).
		Str("account", util.GetAccount(ctx))
}

func Info(ctx context.Context) *zerolog.Event {
	return fillTraceInfos(ctx, defaultLogger.Info())
}

func Error(ctx context.Context) *zerolog.Event {
	return fillTraceInfos(ctx, defaultLogger.Error())
}

func Debug(ctx context.Context) *zerolog.Event {
	return fillTraceInfos(ctx, defaultLogger.Debug())
}

func Warn(ctx context.Context) *zerolog.Event {
	return fillTraceInfos(ctx, defaultLogger.Warn())
}

// NewHTTPServerLogger create a http server logger
func NewHTTPServerLogger() *log.Logger {
	return log.New(&httpServerLogger{}, "", 0)
}

// NewRedisLogger create a redis logger
func NewRedisLogger() *redisLogger {
	return &redisLogger{}
}

// NewEntLogger create a ent logger
func NewEntLogger() *entLogger {
	return &entLogger{}
}

// URLValues create a url.Values log event
func URLValues(query url.Values) *zerolog.Event {
	if len(query) == 0 {
		return zerolog.Dict()
	}
	return zerolog.Dict().Fields(logMask.URLValues(query))
}

// Struct create a struct log event
func Struct(data interface{}) *zerolog.Event {
	if data == nil {
		return zerolog.Dict()
	}

	// 转换出错忽略
	m, _ := logMask.Struct(data)
	return zerolog.Dict().Fields(m)
}
