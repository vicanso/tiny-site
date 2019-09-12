// Copyright 2019 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"time"

	"github.com/vicanso/elton"

	"github.com/vicanso/go-axios"

	"github.com/vicanso/tiny-site/config"
	"github.com/vicanso/tiny-site/helper"
)

var (
	// LocationIns location ins
	LocationIns *axios.Instance
)

const (
	locationService = "location"
)

// 相关的URL
const (
	locationURL = "/ip-locations/json/:ip"
)

type (
	// Location location
	Location struct {
		IP       string `json:"ip,omitempty"`
		Country  string `json:"country,omitempty"`
		Province string `json:"province,omitempty"`
		City     string `json:"city,omitempty"`
		ISP      string `json:"isp,omitempty"`
	}
)

func init() {
	locationBaseURL := config.GetString("location.baseURL")
	LocationIns = helper.NewInstance(locationService, locationBaseURL, 5*time.Second)
}

// GetLocationByIP get location by ip
func GetLocationByIP(ip string, c *elton.Context) (lo Location, err error) {
	conf := &axios.Config{
		URL: locationURL,
		Params: map[string]string{
			"ip": ip,
		},
	}
	helper.AttachWithContext(conf, c)
	resp, err := LocationIns.Request(conf)
	if err != nil {
		return
	}
	err = resp.JSON(&lo)
	if err != nil {
		return
	}
	return
}
