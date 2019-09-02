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
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"github.com/vicanso/hes"
	"github.com/vicanso/tiny-site/config"
	"github.com/vicanso/tiny-site/util"
)

var (
	fontPath string
)

const (
	captchaKeyPrefix = "captcha-"
)

type (
	// CaptchaInfo captcha info
	CaptchaInfo struct {
		Data []byte `json:"data,omitempty"`
		// json输出时，忽略此字段
		Value string `json:"-"`
		ID    string `json:"id,omitempty"`
		Type  string `json:"type,omitempty"`
	}
)

func init() {
	fontPath = config.GetString("resources.font")
	draw2d.SetFontFolder(fontPath)
}

// createCaptcha create captcha image
func createCaptcha(fontColor, bgColor color.Color, width, height int, text string) (img *image.RGBA, err error) {
	img = image.NewRGBA(image.Rect(0, 0, width, height))
	gc := draw2dimg.NewGraphicContext(img)
	draw2dkit.RoundedRectangle(gc, 0, 0, float64(width), float64(height), 0, 0)
	gc.SetFillColor(bgColor)
	gc.Fill()

	gc.FillStroke()

	// Set the font luximbi.ttf
	gc.SetFontData(draw2d.FontData{Name: "luxi", Family: draw2d.FontFamilyMono, Style: draw2d.FontStyleBold | draw2d.FontStyleItalic})

	gc.SetFillColor(fontColor)
	fontCount := len(text)
	offset := 10
	eachFontWidth := (width - 2*offset) / fontCount
	fontSize := float64(eachFontWidth) * 1.2
	for index, ch := range text {
		newFontSize := float64(rand.Int63n(40)+80) / 100 * fontSize
		gc.SetFontSize(newFontSize)
		angle := float64(rand.Int63n(20))/100 - 0.1
		offsetX := float64(eachFontWidth + index*eachFontWidth + int(rand.Int63n(10)) - 10)
		offsetY := float64(height) + float64(rand.Int63n(10)) - float64(15)
		if offsetY > float64(height) || offsetX < float64(height)-newFontSize {
			offsetY = float64(height)
		}
		gc.Rotate(angle)
		gc.FillStringAt(string(ch), offsetX, offsetY)
	}

	gc.SetStrokeColor(fontColor)
	gc.SetLineWidth(1)
	for index := 0; index < 8; index++ {
		gc.BeginPath() // Initialize a new path
		x1 := float64(rand.Int31n(int32(width / 2)))
		y1 := float64(rand.Int31n(int32(height)))

		x2 := float64(rand.Int31n(int32(width/2)) + int32(width/2))
		y2 := float64(rand.Int31n(int32(height)))
		gc.MoveTo(x1, y1)
		gc.LineTo(x2, y2)
		gc.Close()
		gc.FillStroke()
	}

	return
}

func parseColor(s string) (c color.RGBA, err error) {
	arr := strings.Split(s, ",")
	if len(arr) != 3 {
		err = hes.New("color is invalid")
		return
	}
	c.A = 0xff
	for index, v := range arr {
		value, e := strconv.Atoi(v)
		if e != nil {
			err = hes.Wrap(e)
			return
		}
		if value > 255 || value < 0 {
			err = hes.New("color value is invalid")
			return
		}
		switch index {
		case 0:
			c.R = uint8(value)
		case 1:
			c.G = uint8(value)
		default:
			c.B = uint8(value)
		}
	}
	return
}

// GetCaptcha get captcha
func GetCaptcha(fontColor, bgColor string) (info *CaptchaInfo, err error) {
	value := util.RandomDigit(4)
	fc, err := parseColor(fontColor)
	if err != nil {
		return
	}
	bc, err := parseColor(bgColor)
	if err != nil {
		return
	}

	img, err := createCaptcha(fc, bc, 80, 40, value)
	if err != nil {
		return
	}
	buffer := new(bytes.Buffer)
	err = jpeg.Encode(buffer, img, nil)
	if err != nil {
		return
	}
	id := util.GenUlid()
	err = redisSrv.Set(captchaKeyPrefix+id, value, 2*time.Minute)
	if err != nil {
		return
	}
	info = &CaptchaInfo{
		Data:  buffer.Bytes(),
		Value: value,
		ID:    id,
		Type:  "jpeg",
	}
	return
}

// ValidateCaptcha validate the captch
func ValidateCaptcha(id, value string) (valid bool, err error) {
	// 开发环境允许万能验证码
	if util.IsDevelopment() && value == "1053" {
		return true, nil
	}
	data, err := redisSrv.GetAndDel(captchaKeyPrefix + id)
	if err != nil {
		return
	}
	valid = data == value
	return
}
