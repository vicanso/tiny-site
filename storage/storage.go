// Copyright 2022 tree xie
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

package storage

import (
	"bytes"
	"context"
	"image"
	"net/url"
	"strings"
	"sync"

	"github.com/vicanso/go-axios"
	"github.com/vicanso/hes"
	"github.com/vicanso/tiny-site/ent/storage"
	"github.com/vicanso/tiny-site/helper"
	"github.com/vicanso/tiny-site/log"
	"github.com/vicanso/tiny-site/schema"
	"github.com/vicanso/upstream"
)

var httpUpstreams = sync.Map{}

var finders = sync.Map{}

func newHTTPImageFinder(name, uri string) (ImageFinder, error) {
	urlInfo, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	uh := &upstream.HTTP{
		Ping: urlInfo.Path,
	}
	for _, host := range strings.Split(urlInfo.Host, ",") {
		uh.Add(urlInfo.Scheme + "://" + host)
	}
	httpUpstreams.Store(name, uh)
	uh.OnStatus(func(status int32, upstream *upstream.HTTPUpstream) {
		log.Info(context.Background()).
			Str("addr", upstream.URL.String()).
			Int32("status", status).
			Msg("")
	})
	// 先执行一次health check
	uh.DoHealthCheck()
	go uh.StartHealthCheck()
	return func(ctx context.Context, params ...string) (*Image, error) {
		if len(params) == 0 {
			return nil, hes.New("request uri can not be empty")
		}
		requestURI, err := url.QueryUnescape(params[0])
		if err != nil {
			return nil, err
		}
		u := uh.PolicyRoundRobin()
		if u == nil {
			return nil, hes.New("get http upstream fail")
		}
		return GetImageFromURL(ctx, u.URL.String()+requestURI)
	}, nil
}
func GetImageFromURL(ctx context.Context, url string) (*Image, error) {
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
	return &Image{
		Type:   t,
		Size:   len(resp.Data),
		Width:  img.Bounds().Dx(),
		Height: img.Bounds().Dy(),
		Data:   resp.Data,
		img:    img,
	}, nil
}

func InitImageFinder(ctx context.Context) error {
	result, err := helper.EntGetClient().Storage.Query().
		Where(storage.StatusEQ(schema.StatusEnabled)).
		All(ctx)
	if err != nil {
		return err
	}
	for _, item := range result {
		var finder ImageFinder
		var err error
		switch item.Category {
		case "http":
			finder, err = newHTTPImageFinder(item.Name, item.URI)
		}
		// 初始化finder失败时，只输出日志
		if err != nil {
			log.Error(ctx).
				Str("name", item.Name).
				Str("uri", item.URI).
				Msg("init finder fail")
			continue
		}
		finders.Store(item.Name, finder)
	}
	return nil
}

func GetFinder(name string) (ImageFinder, error) {
	value, ok := finders.Load(name)
	if !ok {
		return nil, hes.New("finder is not found")
	}
	fn, ok := value.(ImageFinder)
	if !ok {
		return nil, hes.New("finder is invalid")
	}
	return fn, nil
}
