package middleware

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vicanso/tiny-site/util"
)

func TestNewJSONParser(t *testing.T) {
	conf := JSONParserConfig{
		Limit: 10 * 1024,
	}
	fn := NewJSONParser(conf)

	t.Run("post body data parse", func(t *testing.T) {
		body := []byte(`{"account": "vicanso"}`)
		r := httptest.NewRequest(http.MethodPost, "http://aslant.site/", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, r)
		fn(ctx)
		data := util.GetRequestBody(ctx)
		if !bytes.Equal(body, data) {
			t.Fatalf("json parser fail")
		}
	})

	t.Run("read post body fail", func(t *testing.T) {
		reader := util.NewErrorReadCloser(errors.New("read error"))
		r := httptest.NewRequest(http.MethodPost, "http://aslant.site/", reader)
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, r)
		fn(ctx)
		if ctx.GetStatusCode() != http.StatusInternalServerError {
			t.Fatalf("read error's status code should be 500")
		}
	})

	t.Run("post nil", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "http://aslant.site/", nil)
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, r)
		fn(ctx)
		data := util.GetRequestBody(ctx)
		if len(data) != 0 {
			t.Fatalf("the request body should be empty")
		}
	})

	t.Run("post body data parse over limit", func(t *testing.T) {
		conf := JSONParserConfig{
			Limit: 10,
		}
		limitFn := NewJSONParser(conf)
		body := []byte(`{"account": "vicanso"}`)
		r := httptest.NewRequest(http.MethodPost, "http://aslant.site/", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, r)
		limitFn(ctx)
		if ctx.GetStatusCode() != http.StatusRequestEntityTooLarge {
			t.Fatalf("request post data should be too large")
		}
	})

	t.Run("post request not json should pass", func(t *testing.T) {
		body := []byte(`{"account": "vicanso"}`)
		r := httptest.NewRequest(http.MethodPost, "http://aslant.site/", bytes.NewReader(body))
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, r)
		fn(ctx)
		if util.GetRequestBody(ctx) != nil {
			t.Fatalf("post data not json, json parser should be passed")
		}
	})

	t.Run("get request should pass", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://aslant.site/users/v1/me", nil)
		ctx := util.NewContext(nil, r)
		fn(ctx)
		if util.GetRequestBody(ctx) != nil {
			t.Fatalf("json parser should be passed")
		}
	})

}
