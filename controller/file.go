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
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"path/filepath"

	"github.com/vicanso/elton"
	"github.com/vicanso/tiny-site/router"
)

type (
	fileCtrl struct{}
	fileInfo struct {
		Data   []byte `json:"data,omitempty"`
		Type   string `json:"type,omitempty"`
		Size   int    `json:"size,omitempty"`
		Width  int    `json:"width,omitempty"`
		Height int    `json:"height,omitempty"`
	}
)

func init() {
	ctrl := fileCtrl{}
	g := router.NewGroup("/files")
	g.POST("/v1/upload", ctrl.upload)

	g.POST("/v1/zones", ctrl.createZone)

}

func (ctrl fileCtrl) upload(c *elton.Context) (err error) {
	file, header, err := c.Request.FormFile("filename")
	if err != nil {
		return
	}
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	t := filepath.Ext(header.Filename)
	if t != "" {
		t = t[1:]
	}
	info := &fileInfo{
		Data: buf,
		Type: t,
		Size: len(buf),
	}

	r := bytes.NewBuffer(buf)
	var img image.Image
	switch t {
	case "png":
		img, err = png.Decode(r)
	case "jpeg":
		img, err = jpeg.Decode(r)
	}
	if err != nil {
		return
	}
	if img != nil {
		info.Width = img.Bounds().Dx()
		info.Height = img.Bounds().Dy()
	}
	c.Body = info
	return
}

func (ctrl fileCtrl) createZone(c *elton.Context) (err error) {
	return
}
