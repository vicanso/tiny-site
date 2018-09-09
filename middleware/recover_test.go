package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kataras/iris"

	"github.com/vicanso/tiny-site/util"
)

type (
	recoverErrorStruct struct {
		Exception bool
		Message   string
		Stack     []string
	}
)

func TestNewRecover(t *testing.T) {
	t.Run("panic error", func(t *testing.T) {
		fn := NewRecover()
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, nil)
		message := "abcd"
		ctx.AddHandler(func(ctx iris.Context) {
			ctx.Next()
		}, fn, func(ctx iris.Context) {
			panic(errors.New(message))
		})
		ctx.Next()
		if ctx.GetStatusCode() != http.StatusInternalServerError {
			t.Fatalf("the status code should be 500")
		}
		m := &recoverErrorStruct{}
		err := json.Unmarshal(w.Body.Bytes(), m)
		if err != nil {
			t.Fatalf("the response is not json, %v", err)
		}
		if !m.Exception || m.Message != message || len(m.Stack) == 0 {
			t.Fatalf("the exception response is wrong")
		}
	})

	t.Run("panic nil", func(t *testing.T) {
		fn := NewRecover()
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, nil)
		ctx.AddHandler(func(ctx iris.Context) {
			ctx.Next()
		}, fn, func(ctx iris.Context) {
			panic(nil)
		})
		ctx.Next()
		if ctx.GetStatusCode() != http.StatusOK {
			t.Fatalf("status should be 200")
		}
		if util.GetBody(ctx) != nil {
			t.Fatalf("body should be nil")
		}
	})

	t.Run("is stop", func(t *testing.T) {
		fn := NewRecover()
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, nil)
		message := "abcd"
		ctx.AddHandler(func(ctx iris.Context) {
			ctx.Next()
		}, fn, func(ctx iris.Context) {
			ctx.StopExecution()
			panic(errors.New(message))
		})
		ctx.Next()
		if ctx.GetStatusCode() != http.StatusOK {
			t.Fatalf("status should be 200")
		}
		if util.GetBody(ctx) != nil {
			t.Fatalf("body should be nil")
		}
	})

	t.Run("panic not error", func(t *testing.T) {
		fn := NewRecover()
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, nil)
		ctx.AddHandler(func(ctx iris.Context) {
			ctx.Next()
		}, fn, func(ctx iris.Context) {
			panic(map[string]string{
				"a": "1",
			})
		})
		ctx.Next()

		if ctx.GetStatusCode() != http.StatusInternalServerError {
			t.Fatalf("the status code should be 500")
		}
		m := &recoverErrorStruct{}
		err := json.Unmarshal(w.Body.Bytes(), m)
		if err != nil {
			t.Fatalf("the response is not json, %v", err)
		}
		if !m.Exception ||
			m.Message != "map[a:1]" ||
			len(m.Stack) == 0 {
			t.Fatalf("the exception response is wrong")
		}
	})

}
