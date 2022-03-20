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
	"io"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/vicanso/go-axios"
	"github.com/vicanso/hes"
	"github.com/vicanso/tiny-site/ent/storage"
	"github.com/vicanso/tiny-site/helper"
	"github.com/vicanso/tiny-site/log"
	"github.com/vicanso/tiny-site/schema"
	"github.com/vicanso/upstream"

	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

// 记录所有的http upstream
var httpUpstreams = sync.Map{}

// 记录所有的mongodb client
var mongoClients = sync.Map{}

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

func newMinioImageFinder(_, uri string) (ImageFinder, error) {
	urlInfo, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	accessKey := urlInfo.Query().Get("accessKey")
	secretKey := urlInfo.Query().Get("secretKey")
	minioClient, err := minio.New(urlInfo.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	return func(ctx context.Context, params ...string) (*Image, error) {
		if len(params) != 2 {
			return nil, hes.New("minio params is invalid")
		}
		obj, err := minioClient.GetObject(ctx, params[0], params[1], minio.GetObjectOptions{})
		if err != nil {
			return nil, err
		}
		buf, err := io.ReadAll(obj)
		if err != nil {
			return nil, err
		}
		return NewImageFromBytes(buf)
	}, nil
}

func newMongoImageFinder(name, uri string) (ImageFinder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cs, err := connstring.ParseAndValidate(uri)
	if err != nil {
		return nil, err
	}
	if len(cs.Database) == 0 {
		return nil, hes.New("database can not be nil")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	mongoClients.Store(name, client)
	return func(ctx context.Context, params ...string) (*Image, error) {
		if len(params) == 0 {
			return nil, hes.New("gridfs params is invalid")
		}
		db := client.Database(cs.Database)
		collection := options.DefaultName
		if len(params) > 1 {
			collection = params[1]
		}
		bucket, err := gridfs.NewBucket(db, options.GridFSBucket().SetName(collection))
		if err != nil {
			return nil, err
		}
		buffer := bytes.Buffer{}
		id, err := primitive.ObjectIDFromHex(params[0])
		if err != nil {
			return nil, err
		}
		_, err = bucket.DownloadToStream(id, &buffer)
		if err != nil {
			return nil, err
		}
		return NewImageFromBytes(buffer.Bytes())
	}, nil
}

func newOSSImageFinder(_, uri string) (ImageFinder, error) {
	urlInfo, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	accessKey := urlInfo.Query().Get("accessKey")
	secretKey := urlInfo.Query().Get("secretKey")
	if len(accessKey) == 0 || len(secretKey) == 0 {
		return nil, hes.New("access key and secret key can not be nil")
	}
	client, err := oss.New(urlInfo.Hostname(), accessKey, secretKey)
	if err != nil {
		return nil, err
	}

	return func(_ context.Context, params ...string) (*Image, error) {
		if len(params) != 2 {
			return nil, hes.New("oss params is invalid")
		}
		bucket, err := client.Bucket(params[0])
		if err != nil {
			return nil, err
		}
		r, err := bucket.GetObject(params[1])
		if err != nil {
			return nil, err
		}
		buf, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}

		return NewImageFromBytes(buf)
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
	return NewImageFromBytes(resp.Data)
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
		// 如果以$开头，则从env中获取
		if strings.HasPrefix(item.URI, "$") {
			item.URI = os.Getenv(item.URI[1:])
		}
		var err error
		switch item.Category {
		case schema.StorageCategoryHTTP:
			finder, err = newHTTPImageFinder(item.Name, item.URI)
		case schema.StorageCategoryMinio:
			finder, err = newMinioImageFinder(item.Name, item.URI)
		case schema.StorageCategoryGridfs:
			finder, err = newMongoImageFinder(item.Name, item.URI)
		case schema.StorageCategoryOSS:
			finder, err = newOSSImageFinder(item.Name, item.URI)
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
