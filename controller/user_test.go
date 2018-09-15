package controller

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kataras/iris"

	"github.com/vicanso/session"
	"github.com/vicanso/tiny-site/config"
	"github.com/vicanso/tiny-site/model"
	"github.com/vicanso/tiny-site/util"
)

func TestUserCtrl(t *testing.T) {
	ctrl := userCtrl{}
	cookies := []string{}
	account := util.RandomString(10)
	password := util.Sha1(config.GetString("app") + "123456")
	t.Run("getInfo", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/users/v1/me", nil)
		w := httptest.NewRecorder()
		sess := &session.Session{}
		ctx := util.NewContext(w, r)
		util.SetSession(ctx, sess)
		ctrl.getInfo(ctx)

		cookies = ctx.ResponseWriter().Header()["Set-Cookie"]
		userInfo := util.GetBody(ctx).(*userInfoResponse)
		if !userInfo.Anonymous {
			t.Fatalf("user info should be anonymous")
		}

		if userInfo.Date == "" {
			t.Fatalf("user info's date should not be empty")
		}
	})

	t.Run("getAvatar", func(t *testing.T) {
		ctx := util.NewResContext()
		ctrl.getAvatar(ctx)
		if !strings.HasPrefix(ctx.GetContentType(), "image/jpeg") {
			t.Fatalf("the content type should be jpeg")
		}
		buf := util.GetBody(ctx).([]byte)

		if base64.StdEncoding.EncodeToString(buf) != avatar {
			t.Fatalf("the content is wrong")
		}
	})

	t.Run("getLoginToken", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/users/v1/me/token", nil)
		w := httptest.NewRecorder()
		sess := &session.Session{}
		ctx := util.NewContext(w, r)
		util.SetSession(ctx, sess)
		ctrl.getLoginToken(ctx)
		data := util.GetBody(ctx).(iris.Map)
		if len(data["token"].(string)) != 8 {
			t.Fatalf("get login token fail")
		}
	})

	t.Run("register", func(t *testing.T) {
		m := map[string]string{
			"account":  account,
			"password": password,
		}
		buf, _ := json.Marshal(m)
		r := httptest.NewRequest(http.MethodPost, "http://127.0.0.1/users/v1/me", nil)
		w := httptest.NewRecorder()
		sess := session.Mock(session.M{
			"fetched": true,
			"data":    session.M{},
		})

		ctx := util.NewContext(w, r)
		util.SetRequestBody(ctx, []byte("{}"))
		util.SetSession(ctx, sess)
		ctrl.register(ctx)
		if ctx.GetStatusCode() != http.StatusBadRequest {
			t.Fatalf("use bad params, register should fail")
		}

		ctx = util.NewContext(w, r)
		util.SetRequestBody(ctx, buf)
		util.SetSession(ctx, sess)
		ctrl.register(ctx)
		if ctx.GetStatusCode() != http.StatusCreated {
			t.Fatalf("register fail, %v", util.GetBody(ctx))
		}
	})

	t.Run("login param invalid", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "http://127.0.0.1/users/v1/me/login", nil)
		w := httptest.NewRecorder()

		ctx := util.NewContext(w, r)
		util.SetRequestBody(ctx, []byte("{}"))
		ctrl.doLogin(ctx)
		errData := util.GetBody(ctx).(iris.Map)
		if ctx.GetStatusCode() != http.StatusBadRequest || errData["category"] != util.ErrCategoryValidate {
			t.Fatalf("login params should be invalid")
		}
	})

	t.Run("login not track cookie", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "http://127.0.0.1/users/v1/me/login", nil)
		w := httptest.NewRecorder()

		ctx := util.NewContext(w, r)
		// 无track cookie的出错
		util.SetRequestBody(ctx, []byte(`{
			"account": "vicanso",
			"password": "12341234"
		}`))
		ctrl.doLogin(ctx)
		errData := util.GetBody(ctx).(iris.Map)
		if ctx.GetStatusCode() != http.StatusBadRequest ||
			errData["message"] != "track key is not found" {
			t.Fatalf("no track key should return error")
		}
	})

	t.Run("login token is nil", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "http://127.0.0.1/users/v1/me/login", nil)
		w := httptest.NewRecorder()

		ctx := util.NewContext(w, r)
		sess := session.Mock(session.M{
			"fetched": true,
			"data":    session.M{},
		})
		util.SetRequestBody(ctx, []byte(`{
			"account": "vicanso",
			"password": "12341234"
		}`))
		util.SetSession(ctx, sess)
		for _, v := range cookies {
			arr := strings.Split(v, ";")
			arr = strings.Split(arr[0], "=")
			r.AddCookie(&http.Cookie{
				Name:  arr[0],
				Value: arr[1],
			})
		}

		ctrl.doLogin(ctx)
		errData := util.GetBody(ctx).(iris.Map)
		if ctx.GetStatusCode() != http.StatusBadRequest ||
			errData["message"] != "login token can not be nil" {
			t.Fatalf("no login token should return error")
		}
	})

	t.Run("login account/password is wrong", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "http://127.0.0.1/users/v1/me/login", nil)
		w := httptest.NewRecorder()

		ctx := util.NewContext(w, r)
		sess := session.Mock(session.M{
			"fetched": true,
			"data":    session.M{},
		})
		util.SetRequestBody(ctx, []byte(`{
			"account": "xxxxxxxx",
			"password": "12341234"
		}`))
		util.SetSession(ctx, sess)
		sess.Set(loginTokenKey, util.RandomString(8))
		for _, v := range cookies {
			arr := strings.Split(v, ";")
			arr = strings.Split(arr[0], "=")
			r.AddCookie(&http.Cookie{
				Name:  arr[0],
				Value: arr[1],
			})
		}

		ctrl.doLogin(ctx)
		errData := util.GetBody(ctx).(iris.Map)
		if ctx.GetStatusCode() != http.StatusBadRequest ||
			errData["message"] != "account or password is wrong" {
			t.Fatalf("login not exists account should return error")
		}

		util.SetRequestBody(ctx, []byte(`{
			"account": "`+account+`",
			"password": "12341234"
		}`))

		ctrl.doLogin(ctx)
		errData = util.GetBody(ctx).(iris.Map)
		if ctx.GetStatusCode() != http.StatusBadRequest ||
			errData["message"] != "account or password is wrong" {
			t.Fatalf("login not exists account should return error")
		}
	})

	t.Run("login success then logout", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "http://127.0.0.1/users/v1/me/login", nil)
		w := httptest.NewRecorder()

		ctx := util.NewContext(w, r)
		sess := session.Mock(session.M{
			"fetched": true,
			"data":    session.M{},
		})
		token := util.RandomString(8)
		util.SetRequestBody(ctx, []byte(`{
			"account": "`+account+`",
			"password": "`+util.Sha256(token+password)+`"
		}`))
		util.SetSession(ctx, sess)
		sess.Set(loginTokenKey, token)
		for _, v := range cookies {
			arr := strings.Split(v, ";")
			arr = strings.Split(arr[0], "=")
			r.AddCookie(&http.Cookie{
				Name:  arr[0],
				Value: arr[1],
			})
		}

		ctrl.doLogin(ctx)

		data := util.GetBody(ctx).(*userInfoResponse)
		if data.Account != account {
			t.Fatalf("login fail")
		}

		ctrl.doLogout(ctx)
		data = util.GetBody(ctx).(*userInfoResponse)
		if !data.Anonymous {
			t.Fatalf("logout fail")
		}
	})

	t.Run("refresh", func(t *testing.T) {
		ctx := util.NewResContext()
		sess := session.Mock(session.M{
			"fetched": true,
			"data":    session.M{},
		})
		util.SetSession(ctx, sess)
		ctrl.refresh(ctx)
		if ctx.GetStatusCode() != http.StatusNoContent {
			t.Fatalf("http status should be 204")
		}
		if sess.GetUpdatedAt() == "" {
			t.Fatalf("the updated at should not be empty")
		}
	})

	t.Run("update roles", func(t *testing.T) {
		ctx := util.NewResContext()
		ctx.Params().Set("account", account)
		util.SetRequestBody(ctx, []byte(`{
			"role": "admin",
			"type": "add"
		}`))
		ctrl.updateRoles(ctx)
		if ctx.GetStatusCode() != http.StatusNoContent {
			t.Fatalf("update role fail")
		}
		u := model.User{
			Account: account,
		}
		err := u.First()
		if err != nil {
			t.Fatalf("get user fail, %v", err)
		}
		roles := strings.Join(u.Roles, ",")
		if !strings.Contains(roles, "admin") {
			t.Fatalf("add roles fail")
		}
		util.SetRequestBody(ctx, []byte(`{
			"role": "admin",
			"type": "remove"
		}`))
		ctrl.updateRoles(ctx)
		if ctx.GetStatusCode() != http.StatusNoContent {
			t.Fatalf("update role fail")
		}
		u = model.User{
			Account: account,
		}
		err = u.First()
		if err != nil {
			t.Fatalf("get user fail, %v", err)
		}
		roles = strings.Join(u.Roles, ",")
		if strings.Contains(roles, "admin") {
			t.Fatalf("remove role fail")
		}
	})
}
