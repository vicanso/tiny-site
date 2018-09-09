package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kataras/iris"
	"github.com/vicanso/tiny-site/util"
)

func TestNewStats(t *testing.T) {
	done := false
	conf := StatsConfig{
		OnStats: func(stats *StatsInfo) {
			done = true
			if stats.URI != "http://aslant.site/" ||
				stats.StatusCode != http.StatusOK ||
				stats.Consuming < 100 ||
				stats.Type != 2 ||
				stats.IP == "" {
				t.Fatalf("stats info is wrong")
			}
		},
	}
	fn := NewStats(conf)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://aslant.site/", nil)
	ctx := util.NewContext(w, r)
	message := "success"
	ctx.AddHandler(func(ctx iris.Context) {
		ctx.Next()
	}, fn, func(ctx iris.Context) {
		time.Sleep(time.Millisecond * 100)
		util.Res(ctx, message)
	})
	ctx.Next()
	if !done {
		t.Fatalf("the on stats function isn't called")
	}
}
