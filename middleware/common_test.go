package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vicanso/session"

	"github.com/kataras/iris"

	"github.com/vicanso/tiny-site/util"
)

func TestIsLogined(t *testing.T) {
	t.Run("not logined", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/", nil)
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, r)
		IsLogined(ctx)
		if ctx.GetStatusCode() != http.StatusUnauthorized {
			t.Fatalf("the status code should be 401")
		}
		err := util.GetBody(ctx).(iris.Map)
		if err["category"].(string) != util.ErrCategoryLogic ||
			err["code"].(string) != util.ErrCodeUser {
			t.Fatalf("the http error is not wrong")
		}
	})

	t.Run("logined", func(t *testing.T) {
		sess := session.Mock(session.M{
			"fetched": true,
			"data": session.M{
				"account": "vicanso",
			},
		})
		r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/", nil)
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, r)
		util.SetSession(ctx, sess)
		IsLogined(ctx)
		if ctx.GetStatusCode() != http.StatusOK {
			t.Fatalf("is login check fail")
		}
	})
}

func TestIsAnonymous(t *testing.T) {
	t.Run("not login", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/", nil)
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, r)
		IsAnonymous(ctx)
		if ctx.GetStatusCode() != http.StatusOK {
			t.Fatalf("should be pass to next")
		}
	})

	t.Run("logined", func(t *testing.T) {
		sess := session.Mock(session.M{
			"fetched": true,
			"data": session.M{
				"account": "vicanso",
			},
		})
		r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/", nil)
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, r)
		util.SetSession(ctx, sess)
		IsAnonymous(ctx)
		err := util.GetBody(ctx).(iris.Map)
		if err["category"].(string) != util.ErrCategoryLogic ||
			err["code"].(string) != util.ErrCodeUser {
			t.Fatalf("login should return error")
		}
	})
}

func TestIsSu(t *testing.T) {
	t.Run("not login", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/", nil)
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, r)
		IsSu(ctx)
		if ctx.GetStatusCode() != http.StatusUnauthorized {
			t.Fatalf("need to login")
		}
	})
	t.Run("not su", func(t *testing.T) {
		sess := session.Mock(session.M{
			"fetched": true,
			"data": session.M{
				"account": "vicanso",
			},
		})
		r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/", nil)
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, r)
		util.SetSession(ctx, sess)
		IsSu(ctx)
		if ctx.GetStatusCode() != http.StatusForbidden {
			t.Fatalf("not su should be forbidden")
		}
	})

	t.Run("su", func(t *testing.T) {
		sess := session.Mock(session.M{
			"fetched": true,
			"data": session.M{
				"account": "vicanso",
				"roles":   []string{"su"},
			},
		})
		r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/", nil)
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, r)
		util.SetSession(ctx, sess)
		IsSu(ctx)
		if ctx.GetStatusCode() != http.StatusOK {
			t.Fatalf("su should go to next")
		}
	})
}

func TestWaitFor(t *testing.T) {
	fn := WaitFor(time.Second)
	ctx := util.NewResContext()
	ctx.AddHandler(func(ctx iris.Context) {
		ctx.Next()
	}, fn, func(ctx iris.Context) {
	})
	start := time.Now()
	ctx.Next()
	if time.Now().UnixNano()-start.UnixNano() < time.Second.Nanoseconds() {
		t.Fatalf("wait for middleware fail")
	}
}

func TestIsNilQuery(t *testing.T) {
	t.Run("is nil", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/", nil)
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, r)
		IsNilQuery(ctx)
		if ctx.GetStatusCode() != http.StatusOK {
			t.Fatalf("should pass next when query is nil")
		}
	})

	t.Run("is not nil", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/?a=1", nil)
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, r)
		IsNilQuery(ctx)
		if ctx.GetStatusCode() != http.StatusBadRequest {
			t.Fatalf("should return error when query is not nil")
		}
	})
}
