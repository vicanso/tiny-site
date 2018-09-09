package service

import (
	"testing"

	"github.com/h2non/gock"
)

func TestGetLocationByIP(t *testing.T) {
	defer gock.Off()
	t.Run("get location fail", func(t *testing.T) {
		gock.New("http://ip.taobao.com").
			Get("/service/getIpInfo.php").
			Reply(500).
			BodyString("{}")
		_, err := GetLocationByIP("114.114.114.114")
		if err == nil {
			t.Fatalf("get location should be fail")
		}
	})

	t.Run("ip location not found", func(t *testing.T) {
		gock.New("http://ip.taobao.com").
			Get("/service/getIpInfo.php").
			Reply(200).
			BodyString("")
		_, err := GetLocationByIP("114.114.114.114")
		if err == nil {
			t.Fatalf("ip location not found should return error")
		}
	})

	t.Run("get location success", func(t *testing.T) {
		gock.New("http://ip.taobao.com").
			Get("/service/getIpInfo.php").
			Reply(200).
			BodyString(`{"code":0,"data":{"ip":"114.114.114.114","country":"中国","area":"","region":"江苏","city":"南京","county":"XX","isp":"XX","country_id":"CN","area_id":"","region_id":"320000","city_id":"320100","county_id":"xx","isp_id":"xx"}}`)
		info, err := GetLocationByIP("114.114.114.114")
		if err != nil {
			t.Fatalf("get location by ip fail, %v", err)
		}
		if info.Region != "江苏" {
			t.Fatalf("get location info fail")
		}
	})

}
