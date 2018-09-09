package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kataras/iris"

	"github.com/h2non/gock"
	"github.com/vicanso/tiny-site/util"
)

func TestCommonCtrl(t *testing.T) {
	ctrl := commonCtrl{}
	t.Run("getLocationByIP", func(t *testing.T) {
		defer gock.Off()
		r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/v1/ip-location?ip=abcd", nil)
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, r)
		ctrl.getLocationByIP(ctx)
		data := util.GetBody(ctx).(iris.Map)
		if data["category"] != util.ErrCategoryValidate {
			t.Fatalf("the error category should be validate")
		}
		if ctx.GetStatusCode() != http.StatusBadRequest {
			t.Fatalf("the request query should be invalid")
		}

		gock.New("http://ip.taobao.com").
			Get("/service/getIpInfo.php").
			Reply(500).
			BodyString("{}")
		r = httptest.NewRequest(http.MethodGet, "http://127.0.0.1/v1/ip-location?ip=114.114.114.114", nil)
		w = httptest.NewRecorder()
		ctx = util.NewContext(w, r)
		ctrl.getLocationByIP(ctx)
		data = util.GetBody(ctx).(iris.Map)
		if data["category"] != util.ErrCategoryRequset {
			t.Fatalf("the error category should be request")
		}
		if ctx.GetStatusCode() != http.StatusInternalServerError {
			t.Fatalf("the request response should be 500")
		}

		gock.New("http://ip.taobao.com").
			Get("/service/getIpInfo.php").
			Reply(200).
			BodyString(`{"code":0,"data":{"ip":"114.114.114.114","country":"中国","area":"","region":"江苏","city":"南京","county":"XX","isp":"XX","country_id":"CN","area_id":"","region_id":"320000","city_id":"320100","county_id":"xx","isp_id":"xx"}}`)

		r = httptest.NewRequest(http.MethodGet, "http://127.0.0.1/v1/ip-location?ip=114.114.114.114", nil)
		w = httptest.NewRecorder()
		ctx = util.NewContext(w, r)
		ctrl.getLocationByIP(ctx)
		if ctx.GetStatusCode() != http.StatusOK {
			t.Fatalf("get location fail")
		}
		buf, err := json.Marshal(util.GetBody(ctx))
		if err != nil {
			t.Fatalf("response data is not json, %v", err)
		}

		if string(buf) != `{"ip":"114.114.114.114","country":"中国","region":"江苏","isp":"XX"}` {
			t.Fatalf("respons data is wrong")
		}
	})
}
