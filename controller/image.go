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
	"path/filepath"
	"strconv"
	"strings"

	"github.com/vicanso/elton"
	"github.com/vicanso/hes"
	"github.com/vicanso/tiny-site/router"
	"github.com/vicanso/tiny-site/service"
	"github.com/vicanso/tiny-site/util"
)

type (
	imageCtrl struct{}
)

var (
	supportImageTypes = []string{
		"png",
		"jpeg",
		"jpg",
		"webp",
	}
)
var (
	errImageTypeIsInvalid      = hes.New("image type is invalid")
	errImageTypeIsNotSupported = hes.New("image type isn't supported")
	errImageZoneIsInvalid      = hes.New("image zone is invalid")
)

func init() {
	ctrl := imageCtrl{}
	g := router.NewGroup("/images")

	g.GET("/v1/preview/:fileZoneID/:file", ctrl.preview)
}

func (ctrl imageCtrl) preview(c *elton.Context) (err error) {
	zone, _ := strconv.Atoi(c.Param(fileZoneIDKey))
	file := c.Param("file")
	ext := filepath.Ext(file)
	if ext == "" {
		err = errImageTypeIsInvalid
		return
	}
	imageType := ext[1:]
	// 判断是否支持转换的图片类型
	if !util.ContainsString(supportImageTypes, imageType) {
		err = errImageTypeIsNotSupported
		return
	}

	fileName := strings.Replace(file, ext, "", 1)
	arr := strings.Split(fileName, "-")
	quality := 0
	width := 0
	height := 0
	name := arr[0]
	if len(arr) > 1 {
		quality, err = strconv.Atoi(arr[1])
		if err != nil {
			return
		}
	}
	if len(arr) > 2 {
		width, err = strconv.Atoi(arr[2])
		if err != nil {
			return
		}
	}
	if len(arr) > 3 {
		height, err = strconv.Atoi(arr[3])
		if err != nil {
			return
		}
	}
	// 获取图片数据
	f, err := fileSrv.GetByName(name)
	if err != nil {
		return
	}
	if f.Zone != zone {
		err = errImageZoneIsInvalid
		return
	}
	// 图片转换压缩
	data, err := optimSrv.Image(service.ImageOptimParams{
		Data:       f.Data,
		SourceType: f.Type,
		Type:       imageType,
		Quality:    quality,
		Width:      width,
		Height:     height,
	})
	if err != nil {
		return
	}
	if f.MaxAge != "" {
		c.CacheMaxAge(f.MaxAge)
	}
	c.SetContentTypeByExt(ext)
	c.BodyBuffer = bytes.NewBuffer(data)
	return
}
