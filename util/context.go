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

package util

import (
	"context"

	"github.com/vicanso/elton"
	"github.com/vicanso/tiny-site/config"
)

type contextKey string

const (
	deviceIDKey contextKey = "deviceID"
	traceIDKey  contextKey = "traceID"
	accountKey  contextKey = "account"
)

var sessionConfig = config.MustGetSessionConfig()

// GetTrackID 获取track id
func GetTrackID(c *elton.Context) string {
	trackCookie := sessionConfig.TrackKey
	if trackCookie == "" {
		return ""
	}
	cookie, _ := c.Cookie(trackCookie)
	if cookie == nil {
		return ""
	}
	return cookie.Value
}

// GetSessionID 获取session id
func GetSessionID(c *elton.Context) string {
	cookie, _ := c.Cookie(sessionConfig.Key)
	if cookie == nil {
		return ""
	}
	return cookie.Value
}

func getStringFromContext(ctx context.Context, key contextKey) string {
	v := ctx.Value(key)
	if v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}

// SetDeviceID sets device id to context
func SetDeviceID(ctx context.Context, deviceID string) context.Context {
	return context.WithValue(ctx, deviceIDKey, deviceID)
}

// GetDeviceID gets device is from context
func GetDeviceID(ctx context.Context) string {
	return getStringFromContext(ctx, deviceIDKey)
}

// SetTraceID sets trace id to context
func SetTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

// GetTraceID gets trace id from context
func GetTraceID(ctx context.Context) string {
	return getStringFromContext(ctx, traceIDKey)
}

// SetAccount sets account to context
func SetAccount(ctx context.Context, account string) context.Context {
	return context.WithValue(ctx, accountKey, account)
}

// GetAccount gets account from context
func GetAccount(ctx context.Context) string {
	return getStringFromContext(ctx, accountKey)
}
