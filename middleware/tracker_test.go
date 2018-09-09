package middleware

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/kataras/iris"
	"github.com/vicanso/tiny-site/config"
	"github.com/vicanso/tiny-site/util"
)

func TestNewTracker(t *testing.T) {
	t.Run("response succcess track", func(t *testing.T) {

		done := false
		jsonParser := NewJSONParser(JSONParserConfig{
			Limit: 10 * 1024,
		})
		category := "user-login"
		trackeID := util.RandomString(8)
		fn := NewTracker(category, TrackerConfig{
			Query:    true,
			Params:   true,
			Form:     true,
			Response: true,
			Mask:     regexp.MustCompile(`password`),
			OnTrack: func(info *TrackerInfo) {
				done = true
				if info.Category != category ||
					info.TrackID != trackeID ||
					len(info.Query) != 2 ||
					len(info.Form) != 2 ||
					len(info.Body.(iris.Map)) != 1 {
					t.Fatalf("the tracker data is wrong")
				}
			},
		})
		buf := []byte(`{"account": "vicanso", "password": "mypwd"`)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "http://aslant.site/users/v1/me?type=1&password=mypwd", bytes.NewReader(buf))
		r.Header.Set("Content-Type", "application/json")
		r.AddCookie(&http.Cookie{
			Name:  config.GetTrackKey(),
			Value: trackeID,
		})
		ctx := util.NewContext(w, r)
		ctx.AddHandler(func(ctx iris.Context) {
			ctx.Next()
		}, jsonParser, fn, func(ctx iris.Context) {
			time.Sleep(time.Millisecond * 100)
			util.Res(ctx, iris.Map{
				"account": "vicanso",
			})
		})
		ctx.Next()
		if !done {
			t.Fatalf("tracker middleware fail")
		}
	})

	t.Run("response fail track", func(t *testing.T) {
		done := false
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "http://aslant.site/users/v1/me?type=1&password=mypwd", nil)
		ctx := util.NewContext(w, r)
		onTrack := func(info *TrackerInfo) {
			done = true
			if info.Result != HandleFail {
				t.Fatalf("the tracker data is wrong")
			}
			done = true
		}
		fn := NewDefaultTracker("my-category", onTrack)
		ctx.AddHandler(func(ctx iris.Context) {
			ctx.Next()
		}, fn, func(ctx iris.Context) {
			resErr(ctx, errors.New("abcd"))
		})
		ctx.Next()
		if !done {
			t.Fatalf("tracker middleware fail")
		}
	})
}
