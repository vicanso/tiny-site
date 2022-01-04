// Copyright 2021 tree xie
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

package pipeline

import (
	"bytes"
	"context"
	"errors"
	"image"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/vicanso/go-axios"
	"github.com/vicanso/hes"
	"github.com/vicanso/tiny-site/ent"
)

const (
	PositionTopLeft     = "topLeft"
	PositionTop         = "top"
	PositionTopRight    = "topRight"
	PositionLeft        = "left"
	PositionCenter      = "center"
	PositionRight       = "right"
	PositionBottomLeft  = "bottomLeft"
	PositionBottom      = "bottom"
	PositionBottomRight = "bottomRight"
)

const (
	ImageTypePNG  = "png"
	ImageTypeJPEG = "jpeg"
	ImageTypeWEBP = "webp"
	ImageTypeAVIF = "avif"
)

// 不再执行后续时返回
var ErrAbort = errors.New("abort")

type ImageJob func(context.Context, *ent.Image) (*ent.Image, error)

func Do(ctx context.Context, img *ent.Image, jobs ...ImageJob) (*ent.Image, error) {
	var err error
	for _, fn := range jobs {
		img, err = fn(ctx, img)
		if err != nil {
			// 如果是abort error，则直接返回数据
			if err == ErrAbort {
				return img, nil
			}
			return nil, err
		}
	}
	return img, nil
}

type Parser func([]string, http.Header) (ImageJob, error)

func parseProxy(params []string, _ http.Header) (ImageJob, error) {
	if len(params) != 2 {
		return nil, hes.New("proxy params invalid")
	}
	proxyURL, err := url.QueryUnescape(params[1])
	if err != nil {
		return nil, err
	}
	return NewGetHTTPImage(proxyURL), nil
}

func parseOptim(params []string, _ http.Header) (ImageJob, error) {
	quality := 0
	if len(params) > 1 {
		quality, _ = strconv.Atoi(params[1])
	}
	formats := make([]string, 0)

	if len(params) > 2 {
		formats = append(formats, params[2])
	}
	return NewOptimImage(quality, formats...), nil
}

func parseAutoOptim(params []string, header http.Header) (ImageJob, error) {
	quality := 0
	if len(params) > 1 {
		quality, _ = strconv.Atoi(params[1])
	}
	return NewAutoOptimImage(quality, header), nil
}

func parseFitResize(params []string, _ http.Header) (ImageJob, error) {
	if len(params) < 3 {
		return nil, hes.New("fit resize params invalid")
	}
	width, _ := strconv.Atoi(params[1])
	height, _ := strconv.Atoi(params[2])
	return NewFitResizeImage(width, height), nil
}
func parseFillResize(params []string, _ http.Header) (ImageJob, error) {
	if len(params) < 3 {
		return nil, hes.New("fill resize params invalid")
	}
	width, _ := strconv.Atoi(params[1])
	height, _ := strconv.Atoi(params[2])
	return NewFillResizeImage(width, height), nil
}

func parseBucket(params []string, _ http.Header) (ImageJob, error) {
	if len(params) < 3 {
		return nil, hes.New("bucket params invalid")
	}
	return NewGetEntImage(params[1], params[2]), nil
}

func Parse(tasks []string, header http.Header) ([]ImageJob, error) {
	jobs := make([]ImageJob, 0)
	for _, v := range tasks {
		var fn Parser
		arr := strings.Split(v, "/")
		switch arr[0] {
		case "bucket":
			fn = parseBucket
		case "proxy":
			fn = parseProxy
		case "optim":
			fn = parseOptim
		case "autoOptim":
			fn = parseAutoOptim
		case "fitResize":
			fn = parseFitResize
		case "fillResize":
			fn = parseFillResize
		}
		if fn == nil {
			continue
		}
		job, err := fn(arr, header)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

type imageInfo struct {
	data   []byte
	img    image.Image
	format string
}

func getImageFromURL(ctx context.Context, url string) (*imageInfo, error) {
	resp, err := axios.GetDefaultInstance().GetX(ctx, url)
	if err != nil {
		return nil, err
	}
	if resp.Status != 200 {
		return nil, hes.New("get image fail")
	}
	img, t, err := image.Decode(bytes.NewReader(resp.Data))
	if err != nil {
		return nil, err
	}
	return &imageInfo{
		data:   resp.Data,
		img:    img,
		format: t,
	}, nil
}

func decodeImage(img *ent.Image) (image.Image, error) {
	if len(img.Data) == 0 {
		return nil, hes.New("data of image can not be empty")
	}
	srcImage, _, err := image.Decode(bytes.NewReader(img.Data))
	if err != nil {
		return nil, err
	}
	return srcImage, nil
}

func encodeImage(img image.Image, format string) ([]byte, error) {
	buffer := bytes.Buffer{}
	f := imaging.JPEG
	if format == ImageTypePNG {
		f = imaging.PNG
	}
	err := imaging.Encode(&buffer, img, f)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
