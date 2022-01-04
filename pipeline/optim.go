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

package pipeline

import (
	"context"
	"net/http"
	"strings"

	"github.com/vicanso/tiny-site/config"
	"github.com/vicanso/tiny-site/ent"
	"github.com/vicanso/tiny/pb"
	"google.golang.org/grpc"
)

var tinyConn = mustNewTinyConnection()

func mustNewTinyConnection() *grpc.ClientConn {
	tinyConfig := config.MustGetTinyConfig()
	conn, err := grpc.Dial(tinyConfig.Addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return conn
}

func optim(ctx context.Context, img *ent.Image, quality int, format string) (*ent.Image, error) {
	client := pb.NewOptimClient(tinyConn)
	in := pb.OptimRequest{
		Data:    img.Data,
		Quality: uint32(quality),
	}
	switch format {
	case ImageTypePNG:
		in.Output = pb.Type_PNG
	case ImageTypeWEBP:
		in.Output = pb.Type_WEBP
	case ImageTypeAVIF:
		in.Output = pb.Type_AVIF
	default:
		in.Output = pb.Type_JPEG
	}
	switch img.Type {
	case ImageTypePNG:
		in.Source = pb.Type_PNG
	case ImageTypeWEBP:
		in.Source = pb.Type_WEBP
	case ImageTypeAVIF:
		in.Output = pb.Type_AVIF
	default:
		in.Source = pb.Type_JPEG
	}
	reply, err := client.DoOptim(ctx, &in)
	if err != nil {
		return nil, err
	}
	img.Data = reply.Data
	img.Type = format
	img.Size = len(reply.Data)
	return img, nil
}

func NewAutoOptimImage(quality int, header http.Header) ImageJob {
	return func(ctx context.Context, img *ent.Image) (*ent.Image, error) {
		format := img.Type
		accept := header.Get("Accept")
		acceptWebp := strings.Contains(accept, "image/webp")
		acceptAvif := strings.Contains(accept, "image/avif")

		if acceptAvif {
			format = ImageTypeAVIF
		} else if acceptWebp {
			format = ImageTypeWEBP
			if format == ImageTypePNG {
				quality = 0
			}
		}
		return optim(ctx, img, quality, format)
	}
}

func NewOptimImage(quality int, formats ...string) ImageJob {
	return func(ctx context.Context, img *ent.Image) (*ent.Image, error) {
		format := img.Type
		if len(formats) != 0 {
			format = formats[0]
		}
		return optim(ctx, img, quality, format)
	}
}
