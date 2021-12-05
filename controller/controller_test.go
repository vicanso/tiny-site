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

package controller

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
	se "github.com/vicanso/elton-session"
	"github.com/vicanso/tiny-site/schema"
	"github.com/vicanso/tiny-site/session"
	"github.com/vicanso/hes"
)

func TestListParams(t *testing.T) {
	assert := assert.New(t)
	params := listParams{
		Limit:  10,
		Offset: 100,
		Order:  "-id,name",
		Fields: "id,updatedAt",
	}
	assert.Equal(10, params.GetLimit())
	assert.Equal(100, params.GetOffset())
	assert.Equal([]string{
		"id",
		"updated_at",
	}, params.GetFields())
	assert.Equal(2, len(params.GetOrders()))
}

func newContextAndUserSession() (*elton.Context, *session.UserSession) {
	ctx := context.Background()
	s := se.Session{}
	_, _ = s.Fetch(ctx)
	c := elton.NewContext(nil, nil)
	c.Set(se.Key, &s)
	us := getUserSession(c)
	return c, us
}

func TestIsLogin(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	c, us := newContextAndUserSession()
	assert.False(isLogin(c))
	err := us.SetInfo(ctx, session.UserInfo{
		Account: "treexie",
	})
	assert.Nil(err)
	assert.True(isLogin(c))
}

func TestCheckLogin(t *testing.T) {
	assert := assert.New(t)
	c, us := newContextAndUserSession()
	err := checkLoginMiddleware(c)
	assert.Equal("请先登录", err.(*hes.Error).Message)
	ctx := context.Background()
	err = us.SetInfo(ctx, session.UserInfo{
		Account: "treexie",
	})
	assert.Nil(err)
	done := false
	c.Next = func() error {
		done = true
		return nil
	}
	err = checkLoginMiddleware(c)
	assert.Nil(err)
	assert.True(done)
}

func TestCheckAnonymous(t *testing.T) {
	assert := assert.New(t)
	c, us := newContextAndUserSession()
	done := false
	c.Next = func() error {
		done = true
		return nil
	}
	err := checkAnonymousMiddleware(c)
	assert.Nil(err)
	assert.True(done)
	ctx := context.Background()
	err = us.SetInfo(ctx, session.UserInfo{
		Account: "treexie",
	})
	assert.Nil(err)
	err = checkAnonymousMiddleware(c)
	assert.Equal("已是登录状态，请先退出登录", err.(*hes.Error).Message)
}

func TestNewCheckRolesMiddleware(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	fn := newCheckRolesMiddleware([]string{
		schema.UserRoleAdmin,
	})
	c, us := newContextAndUserSession()
	// 未登录
	err := fn(c)
	assert.Equal("请先登录", err.(*hes.Error).Message)

	// 已登录但无权限
	err = us.SetInfo(ctx, session.UserInfo{
		Account: "treexie",
	})
	assert.Nil(err)
	err = fn(c)
	assert.Equal("禁止使用该功能", err.(*hes.Error).Message)

	// 已登录且权限允许
	err = us.SetInfo(ctx, session.UserInfo{
		Account: "treexie",
		Roles: []string{
			schema.UserRoleAdmin,
		},
	})
	assert.Nil(err)
	done := false
	c.Next = func() error {
		done = true
		return nil
	}
	err = fn(c)
	assert.Nil(err)
	assert.True(done)
}

func TestGetIDFromParams(t *testing.T) {
	assert := assert.New(t)
	c := elton.NewContext(nil, nil)
	c.Params = new(elton.RouteParams)
	c.Params.Add("id", "1")
	id, err := getIDFromParams(c)
	assert.Nil(err)
	assert.Equal(1, id)
}
