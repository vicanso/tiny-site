package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kataras/iris"
	"github.com/vicanso/tiny-site/config"
	"github.com/vicanso/tiny-site/service"
	"github.com/vicanso/tiny-site/util"
)

func TestNewSession(t *testing.T) {

	defaultDuration := time.Hour * 24
	sessConfig := SessionConfig{
		// session cache expires
		Expires: config.GetDurationDefault("session.expires", defaultDuration),
		// the sesion cookie
		Cookie: config.GetSessionCookie(),
		// cookie max age
		CookieMaxAge: config.GetDurationDefault("session.cookie.maxAge", defaultDuration),
		// cookie path
		CookiePath: config.GetCookiePath(),
		// cookie signed keys
		Keys: config.GetSessionKeys(),
	}
	client := service.GetRedisClient()
	t.Run("get session", func(t *testing.T) {
		id := "01CNBNBMNBW92044KPDB8VYKYY"
		buf := []byte(`{
			"a": 1,
			"b": "c"
		}`)
		cmd := client.Set(id, buf, time.Second)
		_, err := cmd.Result()
		if err != nil {
			t.Fatalf("set cache fail, %v", err)
		}
		fn := NewSession(client, sessConfig)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://aslant.site/users/v1/me", nil)
		r.AddCookie(&http.Cookie{
			Name:  "sess",
			Value: id,
		})
		r.AddCookie(&http.Cookie{
			Name:  "sess.sig",
			Value: "rIQ8cMXGRLC22aZeQoU0nZb3BGQ",
		})
		ctx := util.NewContext(w, r)
		ctx.AddHandler(func(ctx iris.Context) {
			ctx.Next()
		}, fn, func(ctx iris.Context) {
			sess := util.GetSession(ctx)
			if sess.GetInt("a") != 1 || sess.GetString("b") != "c" {
				t.Fatalf("get session data fail")
			}
		})
		ctx.Next()
	})

	t.Run("fetch session fail", func(t *testing.T) {
		id := "01CNBNBMNBW92044KPDB8VYKYY"
		buf := []byte(`{
			"a": 1,
			"b": "c",
		}`)
		cmd := client.Set(id, buf, time.Second)
		_, err := cmd.Result()
		if err != nil {
			t.Fatalf("set cache fail, %v", err)
		}
		fn := NewSession(client, sessConfig)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://aslant.site/users/v1/me", nil)
		r.AddCookie(&http.Cookie{
			Name:  "sess",
			Value: id,
		})
		r.AddCookie(&http.Cookie{
			Name:  "sess.sig",
			Value: "rIQ8cMXGRLC22aZeQoU0nZb3BGQ",
		})
		ctx := util.NewContext(w, r)
		ctx.AddHandler(func(ctx iris.Context) {
			ctx.Next()
		}, fn)
		ctx.Next()
		if ctx.GetStatusCode() != http.StatusInternalServerError {
			t.Fatalf("fetch fail should be 500")
		}
		errData := util.GetBody(ctx).(iris.Map)
		if errData["category"].(string) != util.ErrCategorySession {
			t.Fatalf("session error category is wrong")
		}
	})

	t.Run("commit session fail", func(t *testing.T) {
		id := "01CNBNBMNBW92044KPDB8VYKYY"
		buf := []byte(`{
			"a": 1,
			"b": "c"
		}`)
		cmd := client.Set(id, buf, time.Second)
		_, err := cmd.Result()
		if err != nil {
			t.Fatalf("set cache fail, %v", err)
		}
		fn := NewSession(client, sessConfig)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://aslant.site/users/v1/me", nil)
		r.AddCookie(&http.Cookie{
			Name:  "sess",
			Value: id,
		})
		r.AddCookie(&http.Cookie{
			Name:  "sess.sig",
			Value: "rIQ8cMXGRLC22aZeQoU0nZb3BGQ",
		})
		ctx := util.NewContext(w, r)
		ctx.AddHandler(func(ctx iris.Context) {
			ctx.Next()
		}, fn, func(ctx iris.Context) {
			sess := util.GetSession(ctx)
			sess.Set("a", 1)
			client.Close()
		})
		ctx.Next()
		if ctx.GetStatusCode() != http.StatusInternalServerError {
			t.Fatalf("commit fail should be 500 but %d, %v", ctx.GetStatusCode(), util.GetBody(ctx))
		}
		errData := util.GetBody(ctx).(iris.Map)
		if errData["category"].(string) != util.ErrCategorySession {
			t.Fatalf("session error category is wrong")
		}
	})
}
