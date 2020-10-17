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

package controller

import (
	"bytes"
	"io"
	"os"
	"time"

	"github.com/gobuffalo/packr/v2"
	"github.com/vicanso/elton"
	"github.com/vicanso/tiny-site/router"

	eltonMid "github.com/vicanso/elton/middleware"
)

type (
	// assetCtrl asset ctrl
	assetCtrl struct {
	}
	staticFile struct {
		box *packr.Box
	}
)

var (
	box = packr.New("asset", "../web/build")
)

func (sf *staticFile) Exists(file string) bool {
	return sf.box.Has(file)
}
func (sf *staticFile) Get(file string) ([]byte, error) {
	return sf.box.Find(file)
}
func (sf *staticFile) Stat(file string) os.FileInfo {
	return nil
}
func (sf *staticFile) NewReader(file string) (io.Reader, error) {
	buf, err := sf.Get(file)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(buf), nil
}

func init() {
	g := router.NewGroup("")
	ctrl := assetCtrl{}
	g.GET("/", ctrl.index)
	g.GET("/favicon.ico", ctrl.favIcon)

	sf := &staticFile{
		box: box,
	}
	g.GET("/static/*", eltonMid.NewStaticServe(sf, eltonMid.StaticServeConfig{
		Path: "/static",
		// 客户端缓存一年
		MaxAge: 365 * 24 * time.Hour,
		// 缓存服务器缓存一个小时
		SMaxAge:             time.Hour,
		DenyQueryString:     true,
		DisableLastModified: true,
	}))
}

func sendFile(c *elton.Context, file string) (err error) {
	buf, err := box.Find(file)
	if err != nil {
		return
	}
	// 根据文件后续设置类型
	c.SetContentTypeByExt(file)
	c.BodyBuffer = bytes.NewBuffer(buf)
	return
}

func (ctrl assetCtrl) index(c *elton.Context) (err error) {
	c.CacheMaxAge(10 * time.Second)
	return sendFile(c, "index.html")
}

func (ctrl assetCtrl) favIcon(c *elton.Context) (err error) {
	c.CacheMaxAge(time.Hour, 10*time.Minute)
	return sendFile(c, "favicon.ico")
}
