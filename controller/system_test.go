package controller

import (
	"net/http"
	"sort"
	"strings"
	"testing"

	"github.com/kataras/iris"

	"github.com/vicanso/tiny-site/util"
)

func TestSystemCtrl(t *testing.T) {
	ctrl := systemCtrl{}

	t.Run("getStatus", func(t *testing.T) {
		ctx := util.NewResContext()
		ctrl.getStatus(ctx)
		data := util.GetBody(ctx).(iris.Map)
		keys := []string{}
		for key := range data {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		if strings.Join(keys, " ") != "goMaxProcs pid startedAt status uptime version" {
			t.Fatalf("get status fail")
		}
	})

	t.Run("getStats", func(t *testing.T) {
		ctx := util.NewResContext()
		ctrl.getStats(ctx)
		data := util.GetBody(ctx).(iris.Map)
		keys := []string{}
		for key := range data {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		if strings.Join(keys, " ") != "connectingCount heapInuse heapSys routineCount sys" {
			t.Fatalf("get stats fail")
		}
	})

	t.Run("get routes", func(t *testing.T) {
		ctx := util.NewResContext()
		ctrl.getRoutes(ctx)
		data := util.GetBody(ctx).(iris.Map)
		if ctx.GetStatusCode() != http.StatusOK || data["routes"] == nil {
			t.Fatalf("get routes fail")
		}
	})
}
