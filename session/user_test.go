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

package session

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
	se "github.com/vicanso/elton-session"
)

func newUserSession(data string) *UserSession {
	se := se.Session{}
	ctx := context.Background()
	_, _ = se.Fetch(ctx)
	_ = se.Set(ctx, UserSessionInfoKey, data)
	return &UserSession{
		se: &se,
	}
}

func TestUserSession(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	t.Run("get info", func(t *testing.T) {
		us := newUserSession(`{
			"account": "treexie"
		}`)
		info, err := us.GetInfo()
		assert.Nil(err)
		assert.Equal("treexie", info.Account)

		info = us.MustGetInfo()
		assert.Equal("treexie", info.Account)
	})

	t.Run("is logined", func(t *testing.T) {
		us := newUserSession(`{}`)
		assert.False(us.IsLogin())
		us = newUserSession(`{
			"account": "treexie"
		}`)
		assert.True(us.IsLogin())
	})

	t.Run("set info", func(t *testing.T) {
		us := newUserSession(`{}`)
		assert.Equal("", us.MustGetInfo().Account)
		err := us.SetInfo(ctx, UserInfo{
			Account: "treexie",
		})
		assert.Nil(err)
		assert.Equal("treexie", us.MustGetInfo().Account)
	})
}

func TestNewUserSession(t *testing.T) {
	assert := assert.New(t)
	c := elton.NewContext(nil, nil)
	// 未设置session时，user session为空
	us := NewUserSession(c)
	assert.Nil(us)

	c.Set(se.Key, &se.Session{})
	// 读取session并生成user session，并保存至context中
	us = NewUserSession(c)
	assert.NotNil(us)
	assert.NotNil(us.se)

	// 直接从context中读取user session
	us = NewUserSession(c)
	assert.NotNil(us)
	assert.NotNil(us.se)
}
