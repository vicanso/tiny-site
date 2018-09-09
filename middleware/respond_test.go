package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kataras/iris"
	"github.com/vicanso/tiny-site/util"
)

func TestNewRespond(t *testing.T) {
	fn := NewRespond()
	t.Run("response no body", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, nil)
		fn(ctx)
		if w.Code != http.StatusOK || len(w.Body.Bytes()) != 0 {
			t.Fatalf("response no body fail")
		}
	})

	t.Run("response string", func(t *testing.T) {
		text := "abcd"
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, nil)
		util.Res(ctx, text)
		fn(ctx)
		if w.Code != http.StatusOK || text != string(w.Body.Bytes()) {
			t.Fatalf("response string fail")
		}
	})

	t.Run("response bytes", func(t *testing.T) {
		buf := []byte("abcd")
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, nil)
		util.Res(ctx, buf)
		fn(ctx)
		if w.Code != http.StatusOK ||
			!bytes.Equal(buf, w.Body.Bytes()) ||
			w.Header()["Content-Type"][0] != "application/octet-stream" {
			t.Fatalf("response bytes fail")
		}
	})

	t.Run("response json", func(t *testing.T) {
		m := map[string]interface{}{
			"account": "vicanso",
			"age":     18,
			"vip":     true,
		}
		w := httptest.NewRecorder()
		buf := []byte(`{"account":"vicanso","age":18,"vip":true}`)
		ctx := util.NewContext(w, nil)
		util.Res(ctx, m)
		fn(ctx)
		if w.Code != http.StatusOK ||
			!bytes.Equal(buf, w.Body.Bytes()) ||
			w.Header()["Content-Type"][0] != "application/json; charset=UTF-8" {
			t.Fatalf("response json fail")
		}
	})

	t.Run("response error", func(t *testing.T) {
		category := "custom-category"
		message := "error-message"
		code := "custom-code"
		he := &util.HTTPError{
			StatusCode: http.StatusBadRequest,
			Category:   category,
			Message:    message,
			Code:       code,
			Extra: iris.Map{
				"a": 1,
				"b": "2",
			},
		}
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, nil)
		resErr(ctx, he)
		fn(ctx)
		e := &util.HTTPError{}
		err := json.Unmarshal(w.Body.Bytes(), e)
		if err != nil {
			t.Fatalf("response error fail, %v", err)
		}
		if w.Code != http.StatusBadRequest ||
			e.Category != category ||
			e.Message != message ||
			e.Code != code ||
			w.Header()["Content-Type"][0] != "application/json; charset=UTF-8" {
			t.Fatalf("response error fail")
		}
	})
}
