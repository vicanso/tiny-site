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
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/vicanso/elton"
	"github.com/vicanso/hes"
	"github.com/vicanso/tiny-site/config"
	"github.com/vicanso/tiny-site/router"
	"github.com/vicanso/tiny-site/service"
	"github.com/vicanso/tiny-site/util"
)

type (
	imageCtrl struct{}
)

var (
	supportConvertImageTypes = []string{
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
	errImageParamsIsInvalid    = hes.New("image params is invalid")
)

const (
	fileNameKey = "file"
	// 默认的s-maxage为缓存10分钟
	defaultSMaxAge = "10m"
)

type (
	optimImageInfo struct {
		SourceType string `json:"sourceType,omitempty"`
		Type       string `json:"type,omitempty"`
		Quality    int    `json:"quality,omitempty"`
		Width      int    `json:"width,omitempty"`
		Height     int    `json:"height,omitempty"`
		Data       []byte `json:"data,omitempty"`
		MaxAge     string `json:"maxAge,omitempty"`
		Size       int    `json:"size,omitempty"`
	}
)

func init() {
	ctrl := imageCtrl{}
	g := router.NewGroup("/images")

	g.GET("/v1/preview/:"+fileNameKey, ctrl.preview)
	g.GET("/v1/optim/:"+fileNameKey, ctrl.optim)

	g.GET("/v1/config", ctrl.config)
}

func optim(file string) (info *optimImageInfo, err error) {
	ext := filepath.Ext(file)
	if ext == "" {
		err = errImageTypeIsInvalid
		return
	}
	imageType := ext[1:]
	// 判断是否支持转换的图片类型
	if !util.ContainsString(supportConvertImageTypes, imageType) {
		err = errImageTypeIsNotSupported
		return
	}

	fileName := strings.Replace(file, ext, "", 1)
	arr := strings.Split(fileName, "-")
	quality := 0
	width := 0
	height := 0
	crop := 0
	name := arr[0]
	max := 5
	// 参数最多只有5个
	if len(arr) > max {
		err = errImageParamsIsInvalid
		return
	}
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
	if len(arr) > 4 {
		crop, err = strconv.Atoi(arr[4])
		if err != nil {
			return
		}
	}

	// 获取图片数据
	f, err := fileSrv.GetByName(name)
	if err != nil {
		return
	}
	if width != 0 && width > f.Width {
		width = f.Width
	}
	if height != 0 && height > f.Height {
		height = f.Height
	}

	// 图片转换压缩
	data, err := optimSrv.Image(service.ImageOptimParams{
		Data:       f.Data,
		SourceType: f.Type,
		Type:       imageType,
		Quality:    quality,
		Width:      width,
		Height:     height,
		Crop:       crop,
	})
	if err != nil {
		return
	}
	info = &optimImageInfo{
		Data:       data,
		SourceType: f.Type,
		Type:       imageType,
		Quality:    quality,
		Width:      width,
		Height:     height,
		MaxAge:     f.MaxAge,
		Size:       len(data),
	}
	return
}

func (ctrl imageCtrl) preview(c *elton.Context) (err error) {
	file := c.Param(fileNameKey)
	ext := filepath.Ext(file)
	autoDetected := false
	// 如果未设置后缀，则从Accept中判断（不建议使用此方式，尽量由客户端指定后缀）
	// 如果支持webp则webp，否则为jpg
	if ext == "" {
		if strings.Contains(c.GetRequestHeader("Accept"), "image/webp") {
			ext = ".webp"
		} else {
			ext = ".jpeg"
		}
		file += ext
		autoDetected = true
	}
	info, err := optim(file)
	if err != nil {
		return
	}

	if err != nil {
		return
	}
	if info.MaxAge != "" {
		// 如果自动生成后缀的，仅可客户端缓存
		if autoDetected {
			d, _ := time.ParseDuration(info.MaxAge)
			c.SetHeader(elton.HeaderCacheControl, fmt.Sprintf("private, max-age=%d", int(d.Seconds())))
		} else {
			// 设置缓存服务的缓存时长
			c.CacheMaxAge(info.MaxAge, defaultSMaxAge)
		}
	}
	c.SetContentTypeByExt(ext)
	c.BodyBuffer = bytes.NewBuffer(info.Data)
	return
}

func (ctrl imageCtrl) optim(c *elton.Context) (err error) {
	file := c.Param(fileNameKey)
	info, err := optim(file)
	if err != nil {
		return
	}
	c.Body = info
	return
}

func (ctrl imageCtrl) config(c *elton.Context) (err error) {
	c.Body = config.GetImagePreview()
	return
}
