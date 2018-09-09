package middleware

import (
	"net/http"
	"testing"
	"time"

	"github.com/kataras/iris"

	"github.com/vicanso/tiny-site/global"
	"github.com/vicanso/tiny-site/util"
)

func TestNewLimiter(t *testing.T) {
	t.Run("limit pass", func(t *testing.T) {
		fn := NewLimiter(LimiterConfig{
			Max: 1,
		})
		ctx := util.NewContext(nil, nil)
		fn(ctx)
	})

	t.Run("over limit", func(t *testing.T) {
		global.StartApplication()
		fn := NewLimiter(LimiterConfig{
			Max: 0,
		})
		ctx := util.NewResContext()
		fn(ctx)
		if global.IsApplicationRunning() {
			t.Fatalf("the application should be paused after over limit")
		}
		errData := util.GetBody(ctx).(iris.Map)
		if ctx.GetStatusCode() != http.StatusTooManyRequests ||
			errData["message"].(string) != "too many request" {
			t.Fatalf("the respons error should be too many request")
		}
	})
}

func TestResetApplication(t *testing.T) {
	global.PauseApplication()
	if global.IsApplicationRunning() {
		t.Fatalf("application should be pause")
	}
	resetApplicationStatus(time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	if !global.IsApplicationRunning() {
		t.Fatalf("application should resume to running")
	}
}
