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

package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/vicanso/tiny-site/util"
	"github.com/vicanso/hes"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	captchaKeyPrefix = "captcha:"

	errCaptchaCategory = "captcha"
)

type (
	// CaptchaInfo 图形验证码
	CaptchaInfo struct {
		ExpiredAt time.Time `json:"expiredAt"`
		Data      []byte    `json:"data"`
		// json输出时，忽略此字段
		Value string `json:"-"`
		ID    string `json:"id"`
		Type  string `json:"type"`
	}
)

// createCaptcha 创建图形验证码
func createCaptcha(fontColor, bgColor color.Color, width, height int, text string) (image.Image, error) {
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}
	dc := gg.NewContext(width, height)
	// 设置背景色
	dc.SetColor(bgColor)
	dc.Clear()
	fontCount := len(text)
	offset := 10
	eachFontWidth := (width - 2*offset) / fontCount
	fontSize := float64(eachFontWidth) * 1.8
	dc.SetColor(fontColor)
	// 对字符串一个个的填入图片中
	for index, ch := range text {
		// 随机+ -字体大小
		newFontSize := float64(rand.Int63n(40)+80) / 100 * fontSize
		face := truetype.NewFace(font, &truetype.Options{Size: newFontSize})
		dc.SetFontFace(face)
		// 随机偏移角度
		angle := float64(rand.Int63n(20))/100 - 0.1
		// 随机偏移位置
		offsetX := float64(eachFontWidth + index*eachFontWidth + int(rand.Int63n(10)) - 10)
		offsetY := float64(height) + float64(rand.Int63n(10)) - float64(15)
		if offsetY > float64(height) || offsetX < float64(height)-newFontSize {
			offsetY = float64(height)
		}
		dc.Rotate(angle)
		dc.DrawString(string(ch), offsetX, offsetY)
	}
	dc.SetStrokeStyle(gg.NewSolidPattern(fontColor))
	dc.SetLineWidth(1.5)
	// 画线（用于影响机器识别）
	for index := 0; index < 3; index++ {
		x1 := float64(rand.Int31n(int32(width / 2)))
		y1 := float64(rand.Int31n(int32(height)))

		x2 := float64(rand.Int31n(int32(width/2)) + int32(width/2))
		y2 := float64(rand.Int31n(int32(height)))
		dc.DrawLine(x1, y1, x2, y2)
		dc.Stroke()
	}
	return dc.Image(), nil
}

// parseColor 转换颜色
func parseColor(s string) (*color.RGBA, error) {
	arr := strings.Split(s, ",")
	if len(arr) != 3 {
		return nil, hes.New(fmt.Sprintf("非法颜色值，格式必须为：1,1,1，当前为：%s", s), errCaptchaCategory)
	}
	c := &color.RGBA{}
	c.A = 0xff
	for index, v := range arr {
		value, e := strconv.Atoi(strings.TrimSpace(v))
		if e != nil {
			return nil, hes.Wrap(e)
		}
		if value > 255 || value < 0 {
			return nil, hes.New(fmt.Sprintf("非法颜色值，必须>=0 <=255，当前为：%d", value), errCaptchaCategory)
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
	return c, nil
}

// GetCaptcha 获取图形验证码
func GetCaptcha(ctx context.Context, fontColor, bgColor string) (*CaptchaInfo, error) {
	value := util.RandomDigit(4)
	fc, err := parseColor(fontColor)
	if err != nil {
		return nil, err
	}
	bc, err := parseColor(bgColor)
	if err != nil {
		return nil, err
	}

	img, err := createCaptcha(fc, bc, 80, 40, value)
	if err != nil {
		return nil, err
	}
	buffer := new(bytes.Buffer)
	err = jpeg.Encode(buffer, img, nil)
	if err != nil {
		return nil, err
	}
	id := util.GenXID()
	ttl := 5 * time.Minute
	err = redisSrv.Set(ctx, captchaKeyPrefix+id, value, ttl+time.Minute)
	if err != nil {
		return nil, err
	}
	return &CaptchaInfo{
		ExpiredAt: time.Now().Add(ttl),
		Data:      buffer.Bytes(),
		Value:     value,
		ID:        id,
		Type:      "jpeg",
	}, nil
}

// ValidateCaptcha 校验图形验证码是否正确
func ValidateCaptcha(ctx context.Context, id, value string) (bool, error) {
	data, err := redisSrv.GetAndDel(ctx, captchaKeyPrefix+id)
	if err != nil {
		return false, err
	}
	return string(data) == value, nil
}
