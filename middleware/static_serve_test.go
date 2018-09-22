package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kataras/iris"

	"github.com/vicanso/tiny-site/asset"
	"github.com/vicanso/tiny-site/util"
)

func TestStaticServe(t *testing.T) {

	doTest := func(t *testing.T, conf StaticServeConfig) {
		fn := StaticServe(conf)
		r := httptest.NewRequest(http.MethodGet, "http://localhost/?file=index.html", nil)
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, r)
		fn(ctx)
		if ctx.GetStatusCode() != http.StatusOK {
			t.Fatalf("http status code should be 200")
		}
		header := ctx.ResponseWriter().Header()
		if len(header["Cache-Control"]) != 1 ||
			len(header["Etag"]) != 1 ||
			len(header["X-File"]) != 1 {
			t.Fatalf("set http response header fail")
		}
	}
	t.Run("static from asset", func(t *testing.T) {
		conf := StaticServeConfig{
			Asset:   asset.New(),
			MaxAge:  "1h",
			SMaxAge: "1m",
			ETag:    true,
			Header: map[string]string{
				"X-File": "My-Static-File",
			},
		}
		doTest(t, conf)
	})

	t.Run("static from path", func(t *testing.T) {
		conf := StaticServeConfig{
			Path:   "../assets",
			MaxAge: "1h",
			ETag:   true,
			Header: map[string]string{
				"X-File": "My-Static-File",
			},
		}
		doTest(t, conf)
	})

	t.Run("not modified", func(t *testing.T) {
		conf := StaticServeConfig{
			Asset: asset.New(),
			ETag:  true,
		}
		fn := StaticServe(conf)
		r := httptest.NewRequest(http.MethodGet, "http://localhost/?file=index.html", nil)
		r.Header.Set("If-None-Match", "*")
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, r)
		fn(ctx)
		if ctx.GetStatusCode() != http.StatusNotModified {
			t.Fatalf("response should be not modified")
		}
	})

	t.Run("compress", func(t *testing.T) {
		conf := StaticServeConfig{
			Asset:       asset.New(),
			Compression: true,
		}
		fn := StaticServe(conf)
		r := httptest.NewRequest(http.MethodGet, "http://localhost/?file=index.html", nil)
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, r)
		fn(ctx)
		if ctx.GetStatusCode() != http.StatusOK {
			t.Fatalf("gzip response should be 200")
		}
		if ctx.ResponseWriter().Header()["Content-Encoding"][0] != "gzip" {
			t.Fatalf("the content encondig should be gzip")
		}
	})

	t.Run("file name is nil", func(t *testing.T) {
		conf := StaticServeConfig{
			Asset:   asset.New(),
			MaxAge:  "1h",
			SMaxAge: "1m",
		}
		fn := StaticServe(conf)
		r := httptest.NewRequest(http.MethodGet, "http://localhost/", nil)
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, r)
		fn(ctx)
		if ctx.GetStatusCode() != http.StatusBadRequest {
			t.Fatalf("it should return error")
		}
		errData := util.GetBody(ctx).(iris.Map)
		if errData["category"] != util.ErrCategoryLogic ||
			errData["code"] != util.ErrCodeValidate {
			t.Fatalf("error is wrong")
		}
	})

	t.Run("open file fail", func(t *testing.T) {
		conf := StaticServeConfig{
			Path: "../assets",
		}
		fn := StaticServe(conf)
		r := httptest.NewRequest(http.MethodGet, "http://localhost/?file=abc", nil)
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, r)
		fn(ctx)
		if ctx.GetStatusCode() != http.StatusNotFound {
			t.Fatalf("it should return error")
		}
		errData := util.GetBody(ctx).(iris.Map)
		if errData["category"] != util.ErrCategoryLogic ||
			errData["code"] != util.ErrCodeValidate {
			t.Fatalf("error is wrong")
		}
	})

	t.Run("asset file not exists", func(t *testing.T) {
		conf := StaticServeConfig{
			Asset: asset.New(),
		}
		fn := StaticServe(conf)
		r := httptest.NewRequest(http.MethodGet, "http://localhost/?file=abc", nil)
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, r)
		fn(ctx)
		if ctx.GetStatusCode() != http.StatusNotFound {
			t.Fatalf("it should return error")
		}
		errData := util.GetBody(ctx).(iris.Map)
		if errData["category"] != util.ErrCategoryLogic ||
			errData["code"] != util.ErrCodeValidate {
			t.Fatalf("error is wrong")
		}
	})

	t.Run("open default file", func(t *testing.T) {
		conf := StaticServeConfig{
			Path: "../assets",
		}
		fn := StaticServe(conf)
		r := httptest.NewRequest(http.MethodGet, "http://localhost/?file=.", nil)
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, r)
		fn(ctx)
		if ctx.GetStatusCode() != http.StatusOK {
			t.Fatalf("get default fail")
		}
	})

}
