package router

import (
	"net/http"
	"testing"

	"github.com/kataras/iris"
	"github.com/vicanso/tiny-site/global"
	"github.com/vicanso/tiny-site/util"
)

func TestPing(t *testing.T) {
	r := routerList[0]
	if r.Path != "/ping" {
		t.Fatalf("should add ping route for health check")
	}
	global.PauseApplication()
	ctx := util.NewResContext()
	r.Handlers[0](ctx)
	if ctx.GetStatusCode() != http.StatusServiceUnavailable {
		t.Fatalf("should return service unavailable")
	}
	global.StartApplication()
	ctx = util.NewResContext()
	r.Handlers[0](ctx)
	if ctx.GetStatusCode() != http.StatusOK {
		t.Fatalf("should be ok after application running")
	}
}
func TestAdd(t *testing.T) {

	fn := func(ctx iris.Context) {
	}
	testPath := "/test-path"
	Add(http.MethodGet, testPath, fn)
	r := routerList[1]
	if r.Method != http.MethodGet || r.Path != testPath || len(r.Handlers) != 1 {
		t.Fatalf("add router fail")
	}
}

func TestGroup(t *testing.T) {
	isLogin := func(ctx iris.Context) {
	}
	g := NewGroup("/users", isLogin)
	getUserOrders := func(ctx iris.Context) {
	}
	g.Add(http.MethodGet, "/me/orders", getUserOrders)
	r := routerList[2]
	if r.Method != http.MethodGet || r.Path != "/users/me/orders" || len(r.Handlers) != 2 {
		t.Fatalf("add group router fail")
	}
}

func TestList(t *testing.T) {
	if len(List()) != len(routerList) {
		t.Fatalf("list function fail")
	}
}
