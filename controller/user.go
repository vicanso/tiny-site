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

package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/vicanso/hes"
	"github.com/vicanso/tiny-site/middleware"
	"github.com/vicanso/tiny-site/validate"

	"github.com/vicanso/elton"
	"github.com/vicanso/tiny-site/config"
	"github.com/vicanso/tiny-site/cs"
	"github.com/vicanso/tiny-site/router"
	"github.com/vicanso/tiny-site/service"
	"github.com/vicanso/tiny-site/util"
)

type (
	userCtrl struct{}
	// userInfoResp user info response
	userInfoResp struct {
		Anonymous bool     `json:"anonymous,omitempty"`
		Account   string   `json:"account,omitempty"`
		Roles     []string `json:"roles,omitempty"`
		Date      string   `json:"date,omitempty"`
		UpdatedAt string   `json:"updatedAt,omitempty"`
		IP        string   `json:"ip,omitempty"`
		TrackID   string   `json:"trackId,omitempty"`
		LoginAt   string   `json:"loginAt,omitempty"`
	}

	registerUserParams struct {
		Account  string `json:"account,omitempty" valid:"xUserAccount"`
		Password string `json:"password,omitempty" valid:"xUserPassword"`
	}

	loginUserParams struct {
		Account  string `json:"account,omitempty" valid:"xUserAccount"`
		Password string `json:"password,omitempty" valid:"xUserPassword"`
	}

	listUserParams struct {
		Limit   string `json:"limit,omitempty" valid:"xLimit"`
		Keyword string `json:"keyword,omitempty" valid:"xUserAccountKeyword,optional"`
		Role    string `json:"role,omitempty" valid:"xUserRole,optional"`
	}

	updateUserParams struct {
		Roles []string `json:"roles,omitempty" valid:"xUserRoles,optional"`
	}
	listUserLoginRecordParams struct {
		Begin   time.Time `json:"begin,omitempty" valid:"-"`
		End     time.Time `json:"end,omitempty" valid:"-"`
		Account string    `json:"account,omitempty" valid:"xUserAccount,optional"`
		Limit   string    `json:"limit,omitempty" valid:"xLimit"`
		Offset  string    `json:"offset,omitempty" valid:"xOffset"`
	}
)

var (
	errLoginTokenNil = hes.New("login token is nil")
)

func init() {
	g := router.NewGroup("/users", loadUserSession)
	ctrl := userCtrl{}
	// 获取用户列表
	g.GET(
		"/v1",
		shouldBeAdmin,
		ctrl.list,
	)

	// 更新用户信息
	g.PATCH(
		// 因为与/me有冲突，因此路径增加update
		"/v1/update/{userID}",
		shouldBeAdmin,
		ctrl.update,
	)

	// 获取用户信息
	g.GET("/v1/me", ctrl.me)

	// 用户注册
	g.POST(
		"/v1/me",
		newTracker(cs.ActionRegister),
		// 限制相同IP在60秒之内只能调用5次
		newIPLimit(5, 60*time.Second, cs.ActionLogin),
		shouldAnonymous,
		middleware.ValidateCaptch(),
		ctrl.register,
	)
	// 刷新user session的ttl
	g.PATCH(
		"/v1/me",
		ctrl.refresh,
	)

	// 获取登录token
	g.GET(
		"/v1/me/login",
		shouldAnonymous,
		ctrl.getLoginToken,
	)

	// 用户登录
	// 限制3秒只能登录一次（无论成功还是失败）
	loginLimit := newConcurrentLimit([]string{
		"account",
	}, 3*time.Second, cs.ActionLogin)
	g.POST(
		"/v1/me/login",
		middleware.WaitFor(time.Second),
		newTracker(cs.ActionLogin),
		shouldAnonymous,
		loginLimit,
		// 限制相同IP在60秒之内只能调用10次
		newIPLimit(10, 60*time.Second, cs.ActionLogin),
		// 限制10分钟内，相同的账号只允许出错5次
		newErrorLimit(5, 10*time.Minute, func(c *elton.Context) string {
			return standardJSON.Get(c.RequestBody, "account").ToString()
		}),
		middleware.ValidateCaptch(),
		ctrl.login,
	)
	// 用户退出登录
	g.DELETE(
		"/v1/me/logout",
		newTracker(cs.ActionLogout),
		shouldLogined,
		ctrl.logout,
	)

	// 获取客户登录记录
	g.GET(
		"/v1/login-records",
		shouldBeAdmin,
		ctrl.listLoginRecord,
	)
}

// get user info from session
func pickUserInfo(c *elton.Context) (userInfo *userInfoResp) {
	us := getUserSession(c)
	userInfo = &userInfoResp{
		Anonymous: true,
		Date:      now(),
		IP:        c.RealIP(),
		TrackID:   getTrackID(c),
	}
	account := us.GetAccount()
	if account != "" {
		userInfo.Account = account
		userInfo.Roles = us.GetRoles()
		userInfo.Anonymous = false
	}
	return
}

// get user info
func (ctrl userCtrl) me(c *elton.Context) (err error) {
	key := config.GetTrackKey()
	cookie, _ := c.Cookie(key)
	// ulid的长度为26
	if cookie == nil || len(cookie.Value) != 26 {
		uid := util.GenUlid()
		c.AddCookie(&http.Cookie{
			Name:     key,
			Value:    uid,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   365 * 24 * 3600,
		})
		trackRecord := &service.UserTrackRecord{
			UserAgent: c.GetRequestHeader("User-Agent"),
			IP:        c.RealIP(),
			TrackID:   util.GetTrackID(c),
		}
		userSrv.AddTrackRecord(trackRecord)
	}
	c.Body = pickUserInfo(c)
	return
}

// getLoginToken get login token
func (ctrl userCtrl) getLoginToken(c *elton.Context) (err error) {
	us := getUserSession(c)
	// 清除当前session id，确保每次登录的用户都是新的session
	us.ClearSessionID()
	token := util.RandomString(8)
	err = us.SetLoginToken(token)
	if err != nil {
		return
	}
	c.Body = &struct {
		Token string `json:"token,omitempty"`
	}{
		token,
	}
	return
}

func omitUserInfo(u *service.User) {
	u.Password = ""
}

// register user register
func (ctrl userCtrl) register(c *elton.Context) (err error) {
	params := &registerUserParams{}
	err = validate.Do(params, c.RequestBody)
	if err != nil {
		return
	}
	u := &service.User{
		Account:  params.Account,
		Password: params.Password,
	}
	err = userSrv.Add(u)
	if err != nil {
		return
	}
	omitUserInfo(u)
	c.Created(u)
	return
}

// login user login
func (ctrl userCtrl) login(c *elton.Context) (err error) {
	params := &registerUserParams{}
	err = validate.Do(params, c.RequestBody)
	if err != nil {
		return
	}
	us := getUserSession(c)
	token := us.GetLoginToken()
	if token == "" {
		err = errLoginTokenNil
		return
	}
	u, err := userSrv.Login(params.Account, params.Password, token)
	if err != nil {
		return
	}
	loginRecord := &service.UserLoginRecord{
		Account:       params.Account,
		UserAgent:     c.GetRequestHeader("User-Agent"),
		IP:            c.RealIP(),
		TrackID:       util.GetTrackID(c),
		SessionID:     util.GetSessionID(c),
		XForwardedFor: c.GetRequestHeader("X-Forwarded-For"),
	}
	userSrv.AddLoginRecord(loginRecord)
	omitUserInfo(u)
	us.SetAccount(u.Account)
	us.SetRoles(u.Roles)
	c.Body = u
	return
}

// logout user logout
func (ctrl userCtrl) logout(c *elton.Context) (err error) {
	us := getUserSession(c)
	if us != nil {
		err = us.Destroy()
	}
	c.NoContent()
	return
}

// refresh user refresh
func (ctrl userCtrl) refresh(c *elton.Context) (err error) {
	us := getUserSession(c)
	if us == nil {
		c.NoContent()
		return
	}

	scf := config.GetSessionConfig()
	cookie, _ := c.SignedCookie(scf.Key)
	// 如果认证的cookie已过期，则不做刷新
	if cookie == nil {
		c.NoContent()
		return
	}

	err = us.Refresh()
	if err != nil {
		return
	}
	// 更新session
	err = c.AddSignedCookie(&http.Cookie{
		Name:     scf.Key,
		Value:    cookie.Value,
		Path:     scf.CookiePath,
		MaxAge:   int(scf.TTL.Seconds()),
		HttpOnly: true,
	})
	if err != nil {
		return
	}

	c.NoContent()
	return
}

// list user list
func (ctrl userCtrl) list(c *elton.Context) (err error) {
	params := &listUserParams{}
	err = validate.Do(params, c.Query())
	if err != nil {
		return
	}
	limit, _ := strconv.Atoi(params.Limit)
	users, err := userSrv.List(service.UserQueryParams{
		Role:    params.Role,
		Keyword: params.Keyword,
		Limit:   limit,
	})
	if err != nil {
		return
	}
	c.Body = &struct {
		Users []*service.User `json:"users,omitempty"`
	}{
		users,
	}
	return
}

// update user update
func (ctrl userCtrl) update(c *elton.Context) (err error) {
	id, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		return
	}
	params := &updateUserParams{}
	err = validate.Do(params, c.RequestBody)
	if err != nil {
		return
	}
	// 只能su用户才可以添加su权限
	if util.ContainsString(params.Roles, cs.UserRoleSu) {
		roles := getUserSession(c).GetRoles()
		if !util.ContainsString(roles, cs.UserRoleSu) {
			err = hes.New("add su role is forbidden")
			return
		}
	}
	err = userSrv.Update(&service.User{
		ID: uint(id),
	}, map[string]interface{}{
		"roles": params.Roles,
	})
	if err != nil {
		return
	}
	c.NoContent()
	return
}

// listLoginRecord list login record
func (ctrl userCtrl) listLoginRecord(c *elton.Context) (err error) {
	params := &listUserLoginRecordParams{}
	err = validate.Do(params, c.Query())
	if err != nil {
		return
	}
	offset, _ := strconv.Atoi(params.Offset)
	limit, _ := strconv.Atoi(params.Limit)
	query := service.UserLoginRecordQueryParams{
		Account: params.Account,
		Limit:   limit,
		Offset:  offset,
	}
	if !params.Begin.IsZero() {
		query.Begin = util.FormatTime(params.Begin)
	}
	if !params.End.IsZero() {
		query.End = util.FormatTime(params.End)
	}
	result, err := userSrv.ListLoginRecord(query)
	if err != nil {
		return
	}
	count := -1
	if offset == 0 {
		count, err = userSrv.CountLoginRecord(query)
	}
	c.Body = struct {
		Logins []*service.UserLoginRecord `json:"logins,omitempty"`
		Count  int                        `json:"count,omitempty"`
	}{
		result,
		count,
	}
	return
}
