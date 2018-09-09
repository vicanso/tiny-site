package util

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vicanso/tiny-site/config"
	"github.com/vicanso/tiny-site/cs"
	"github.com/vicanso/session"

	"github.com/kataras/iris"
)

type (
	params struct {
		// Account account
		Account string `json:"account" valid:"ascii,runelength(4|10)"`
	}
)

func TestRes(t *testing.T) {
	t.Run("no content", func(t *testing.T) {
		ctx := NewResContext()
		ResNoContent(ctx)
		if ctx.GetStatusCode() != http.StatusNoContent {
			t.Fatalf("no content status code fail")
		}
		if GetBody(ctx) != nil {
			t.Fatalf("no content body should be nil")
		}
	})

	t.Run("response created", func(t *testing.T) {
		ctx := NewResContext()
		m := map[string]string{
			"id": "1",
		}
		ResCreated(ctx, m)
		if ctx.GetStatusCode() != http.StatusCreated {
			t.Fatalf("created status code fail")
		}
		if GetBody(ctx).(map[string]string)["id"] != "1" {
			t.Fatalf("get created response data fail")
		}
	})

	t.Run("ressponse jpeg", func(t *testing.T) {
		ctx := NewResContext()
		buf := []byte("jpeg data")
		ResJPEG(ctx, buf)
		if ctx.GetStatusCode() != http.StatusOK {
			t.Fatalf("response jpeg status code fail")
		}
		if ctx.GetContentType() != "image/jpeg; charset=UTF-8" {
			t.Fatalf("response jpeg content type fail")
		}

		if !bytes.Equal(GetBody(ctx).([]byte), buf) {
			t.Fatalf("response jpeg data fail")
		}
	})

	t.Run("response png", func(t *testing.T) {
		ctx := NewResContext()
		buf := []byte("png data")
		ResPNG(ctx, buf)
		if ctx.GetStatusCode() != http.StatusOK {
			t.Fatalf("response png status code fail")
		}
		if ctx.GetContentType() != "image/png; charset=UTF-8" {
			t.Fatalf("response png content type fail")
		}

		if !bytes.Equal(GetBody(ctx).([]byte), buf) {
			t.Fatalf("response png data fail")
		}
	})

	t.Run("response error", func(t *testing.T) {
		ctx := NewResContext()
		he := &HTTPError{
			Category: "code-category",
			Code:     "custom-error-code",
			Message:  "cus error",
			Extra: iris.Map{
				"uri": "http://aslant.site/",
			},
		}
		ResErr(ctx, he)
		if ctx.GetStatusCode() != http.StatusInternalServerError {
			t.Fatalf("response error status code fail")
		}
		m := GetBody(ctx).(iris.Map)
		if m["code"] != he.Code {
			t.Fatalf("response error data code fail")
		}
		if m["category"] != he.Category {
			t.Fatalf("response error data category fail")
		}
		extra := m["extra"].(iris.Map)
		if extra["uri"].(string) != he.Extra["uri"].(string) {
			t.Fatalf("response error data extra fail")
		}
	})
}

func TestSession(t *testing.T) {
	ctx := NewContext(nil, nil)
	if GetSession(ctx) != nil {
		t.Fatalf("get session should be nil before set")
	}
	sess := &session.Session{}
	SetSession(ctx, sess)
	if GetSession(ctx) != sess {
		t.Fatalf("get/set session fail")
	}
}

func TestGetTrackID(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "http://aslant.site/", nil)
	trackID := "random-string"
	r.AddCookie(&http.Cookie{
		Name:  config.GetTrackKey(),
		Value: trackID,
	})
	ctx := NewContext(nil, r)
	if GetTrackID(ctx) != trackID {
		t.Fatalf("get track id fail")
	}
}

func TestGetAccount(t *testing.T) {
	account := "vicanso"
	ctx := NewResContext()
	if GetAccount(ctx) != "" {
		t.Fatalf("acount should be empty while no session")
	}
	sess := session.Mock(session.M{
		"fetched": true,
		"data": session.M{
			"account": account,
		},
	})
	SetSession(ctx, sess)
	if GetAccount(ctx) != account {
		t.Fatalf("get account fail")
	}
}

func TestRequestBody(t *testing.T) {
	ctx := NewResContext()
	if GetRequestBody(ctx) != nil {
		t.Fatalf("the request body should be nil before set")
	}
	buf := []byte("request body")
	SetRequestBody(ctx, buf)
	if !bytes.Equal(GetRequestBody(ctx), buf) {
		t.Fatalf("get/set request body fail")
	}
}

func TestGetRequestQuery(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "http://aslant.site/", nil)
	ctx := NewContext(nil, r)
	query := GetRequestQuery(ctx)
	if query != nil {
		t.Fatalf("not query should be nil")
	}

	r = httptest.NewRequest(http.MethodGet, "http://aslant.site/?skip=1&limit=10&category=book", nil)
	ctx = NewContext(nil, r)
	query = GetRequestQuery(ctx)
	if len(query) != 3 {
		t.Fatalf("get request query fail")
	}
	if query["skip"] != "1" || query["limit"] != "10" || query["category"] != "book" {
		t.Fatalf("get request query fail")
	}
	// 第二次获取直接从values中获取
	query = GetRequestQuery(ctx)
	if len(query) != 3 {
		t.Fatalf("get request query fail")
	}
	if query["skip"] != "1" || query["limit"] != "10" || query["category"] != "book" {
		t.Fatalf("get request query fail")
	}
}

func TestGetStack(t *testing.T) {
	stack := GetStack(1024)
	if len(stack) == 0 {
		t.Fatalf("get stack fail")
	}
}

func TestSetHeader(t *testing.T) {
	ctx := NewResContext()
	key := "X-Custom-Header"
	SetHeader(ctx, key, "a")
	SetHeader(ctx, key, "b")
	value := ctx.ResponseWriter().Header()[key]
	if len(value) != 1 || value[0] != "b" {
		t.Fatalf("set header fail")
	}
}

func TestRemoveHeader(t *testing.T) {
	ctx := NewResContext()
	key := "X-Custom-Header"
	SetHeader(ctx, key, "a")
	value := ctx.ResponseWriter().Header()[key]
	if len(value) != 1 {
		t.Fatalf("set header fail")
	}
	RemoveHeader(ctx, key)
	value = ctx.ResponseWriter().Header()[key]
	if len(value) != 0 {
		t.Fatalf("remove header fail")
	}
}

func TestSetNoCache(t *testing.T) {
	ctx := NewResContext()
	SetNoCache(ctx)
	value := ctx.ResponseWriter().Header()[cs.HeaderCacheControl][0]
	if value != "no-cache, max-age=0" {
		t.Fatalf("set no cache fail")
	}
}

func TestSetNoStore(t *testing.T) {
	ctx := NewResContext()
	SetNoStore(ctx)
	value := ctx.ResponseWriter().Header()[cs.HeaderCacheControl][0]
	if value != "no-store" {
		t.Fatalf("set no store fail")
	}
}

func TestSetCache(t *testing.T) {
	ctx := NewResContext()
	if err := SetCache(ctx, "abc"); err == nil {
		t.Fatalf("set cache with no duration should return error")
	}
	SetCache(ctx, "10s")
	value := ctx.ResponseWriter().Header()[cs.HeaderCacheControl][0]
	if value != "public, max-age=10" {
		t.Fatalf("set cache fail")
	}
}

func TestSetCacheWithSMaxAge(t *testing.T) {
	ctx := NewResContext()
	if err := SetCacheWithSMaxAge(ctx, "abc", "10s"); err == nil {
		t.Fatalf("set s-maxage cache with no duration should return error")
	}
	if err := SetCacheWithSMaxAge(ctx, "1h", "abc"); err == nil {
		t.Fatalf("set s-maxage cache with no duration should return error")
	}
	SetCacheWithSMaxAge(ctx, "1h", "1m")
	value := ctx.ResponseWriter().Header()[cs.HeaderCacheControl][0]
	if value != "public, max-age=3600, s-maxage=60" {
		t.Fatalf("set s-maxage cache fail")
	}
}
