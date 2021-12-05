// Copyright 2020 tree xie
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

// 此controller提供各静态文件的响应处理，
// 主要是管理系统的前端代码，对于资源等（如图片）尽可能不要打包进入程序

package controller

import (
	"time"

	"github.com/vicanso/elton"
	M "github.com/vicanso/elton/middleware"
	"github.com/vicanso/tiny-site/asset"
	"github.com/vicanso/tiny-site/router"
)

type (
	// assetCtrl asset ctrl
	assetCtrl struct{}
)

var assetFS = M.NewEmbedStaticFS(asset.GetFS(), "dist")

func init() {
	g := router.NewGroup("")
	ctrl := assetCtrl{}
	g.GET("/", ctrl.getIndex)
	g.GET("/favicon.{ext}", ctrl.getFavIcon)

	g.GET("/static/*", M.NewStaticServe(assetFS, M.StaticServeConfig{
		// 客户端缓存一年
		MaxAge: 365 * 24 * time.Hour,
		// 缓存服务器缓存一个小时
		SMaxAge:             time.Hour,
		DenyQueryString:     true,
		DisableLastModified: true,
		EnableStrongETag:    true,
		// 如果静态文件都有版本号，可以指定immutable
		Immutable: true,
	}))
}

// getIndex 首页
func (*assetCtrl) getIndex(c *elton.Context) error {
	err := assetFS.SendFile(c, "index.html")
	if err != nil {
		return err
	}
	c.CacheMaxAge(10 * time.Second)
	return nil
}

// getFavIcon 图标
func (*assetCtrl) getFavIcon(c *elton.Context) error {
	err := assetFS.SendFile(c, "favicon.png")
	if err != nil {
		return err
	}
	c.CacheMaxAge(time.Hour, 10*time.Minute)
	return nil
}
