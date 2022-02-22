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

package controller

import (
	"bytes"
	"context"
	"image"
	"io/ioutil"
	"strings"
	"time"

	"github.com/vicanso/elton"
	"github.com/vicanso/hes"
	"github.com/vicanso/tiny-site/cs"
	"github.com/vicanso/tiny-site/ent"
	"github.com/vicanso/tiny-site/ent/bucket"
	entImage "github.com/vicanso/tiny-site/ent/image"
	"github.com/vicanso/tiny-site/log"
	"github.com/vicanso/tiny-site/pipeline"
	"github.com/vicanso/tiny-site/router"
	"github.com/vicanso/tiny-site/util"
	"github.com/vicanso/tiny-site/validate"
)

type imageCtrl struct{}

type (
	bucketAddParams struct {
		// bucket的名称
		Name string `json:"name" validate:"required,xImageBucket"`
		// 拥有者
		Owners []string `json:"owners" validate:"required,dive,xUserAccount"`
		// 描述
		Description string `json:"description" validate:"required,xImageDescription"`
	}
	bucketUpdateParams struct {
		// 拥有者
		Owners []string `json:"owners" validate:"omitempty,dive,xUserAccount"`
		// 描述
		Description string `json:"description" validate:"omitempty,xImageDescription"`
	}
	bucketListParams struct {
		listParams
	}
	imageAddParams struct {
		Bucket      string `json:"bucket" validate:"required,xImageBucket"`
		Name        string `json:"name" validate:"omitempty,xImageName"`
		Tags        string `json:"tags" validate:"omitempty,xImageTags"`
		Description string `json:"description" validate:"omitempty,xImageDescription"`

		creator string
		data    []byte
	}
	imageListParams struct {
		listParams

		Bucket string `json:"bucket" validate:"required,xImageBucket"`
		Tag    string `json:"tag" validate:"omitempty,xImageTag"`
	}
	imageGetThumbnailParams struct {
		Bucket string `json:"bucket" validate:"required,xImageBucket"`
		// 图片名称
		Name string `json:"name" validate:"omitempty,xImageName"`
		// 缩略图大小
		ThumbnailSize int `json:"thumbnailSize" validate:"omitempty,xImageThumbnailSize" default:"128"`
	}
)

type (
	bucketListResp struct {
		Count   int           `json:"count"`
		Buckets []*ent.Bucket `json:"buckets"`
	}
	imageListResp struct {
		Count  int          `json:"count"`
		Images []*ent.Image `json:"images"`
	}
)

func init() {
	prefix := "/images"
	g := router.NewGroup(prefix, loadUserSession, shouldBeLogin)
	ctrl := imageCtrl{}

	g.GET(
		"/v1/buckets",
		ctrl.listBucket,
	)

	g.POST(
		"/v1/buckets",
		newTrackerMiddleware(cs.ActionBucketAdd),
		ctrl.addBucket,
	)
	g.PATCH(
		"/v1/buckets/{id}",
		newTrackerMiddleware(cs.ActionBucketUpdate),
		ctrl.updateBucket,
	)

	g.POST(
		"/v1",
		newTrackerMiddleware(cs.ActionImageAdd),
		ctrl.addImage,
	)

	g.GET(
		"/v1",
		ctrl.listImage,
	)

	ng := router.NewGroup(prefix)
	ng.GET(
		"/v1/thumbnails/{bucket}/{name}",
		ctrl.getImageThumbnail,
	)
	ng.GET(
		"/v1/pipeline",
		ctrl.pipeline,
	)

}

func (params *bucketListParams) queryAll(ctx context.Context) ([]*ent.Bucket, error) {
	query := getBucketClient().Query()
	query.Limit(params.GetLimit()).
		Offset(params.GetOffset()).
		Order(params.GetOrders()...)
	return query.All(ctx)
}

func (params *bucketListParams) count(ctx context.Context) (int, error) {
	query := getBucketClient().Query()

	return query.Count(ctx)
}

func (params *imageAddParams) save(ctx context.Context) (*ent.Image, error) {
	image, imageType, err := image.Decode(bytes.NewReader(params.data))
	if err != nil {
		return nil, err
	}

	return getImageClient().Create().
		SetBucket(params.Bucket).
		SetName(params.Name).
		SetType(imageType).
		SetSize(len(params.data)).
		SetWidth(image.Bounds().Dx()).
		SetHeight(image.Bounds().Dy()).
		SetTags(params.Tags).
		SetCreator(params.creator).
		SetData(params.data).
		SetDescription(params.Description).
		Save(ctx)
}

func (params *imageListParams) where(query *ent.ImageQuery) *ent.ImageQuery {
	if params.Bucket != "" {
		query.Where(entImage.Bucket(params.Bucket))
	}
	if params.Tag != "" {
		query.Where(entImage.TagsContains(params.Tag))
	}
	return query
}

func (params *imageListParams) queryAll(ctx context.Context) ([]*ent.Image, error) {
	query := getImageClient().Query()
	query = query.Limit(params.GetLimit()).
		Offset(params.GetOffset()).
		Order(params.GetOrders()...)
	params.where(query)
	fields := params.GetFields()
	if len(fields) != 0 {
		result := make([]*ent.Image, 0)
		err := query.Select(fields...).Scan(ctx, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	return query.All(ctx)
}

func (params *imageListParams) count(ctx context.Context) (int, error) {
	query := getImageClient().Query()
	params.where(query)
	return query.Count(ctx)
}

func (params *bucketUpdateParams) updateOneID(ctx context.Context, id int, creator string) (*ent.Bucket, error) {
	result, err := getBucketClient().Query().Where(
		bucket.IDEQ(id),
	).First(ctx)
	if err != nil {
		return nil, err
	}
	if result.Creator != creator {
		return nil, hes.NewWithExcpetion("Forbidden to modify bucket")
	}
	updateOne := getBucketClient().UpdateOneID(id)
	if params.Description != "" {
		updateOne.SetDescription(params.Description)
	}
	if len(params.Owners) != 0 {
		updateOne.SetOwners(params.Owners)
	}
	return updateOne.Save(ctx)
}

func validateBucketForUser(ctx context.Context, bucketName, account string) error {
	if bucketName == "" {
		return hes.New("bucket名称不能为空")
	}
	result, err := getBucketClient().Query().
		Where(bucket.Name(bucketName)).
		First(ctx)
	if err != nil {
		return err
	}
	if len(result.Owners) != 0 && !util.ContainsString(result.Owners, account) {
		return hes.New("无权限添加图片至此bucket")
	}
	return nil
}

func (*imageCtrl) addBucket(c *elton.Context) error {
	params := bucketAddParams{}
	err := validateBody(c, &params)
	if err != nil {
		return err
	}
	account := getUserSession(c).MustGetInfo().Account
	bucket, err := getBucketClient().Create().
		SetName(params.Name).
		SetOwners(params.Owners).
		SetDescription(params.Description).
		SetCreator(account).
		Save(c.Context())
	if err != nil {
		return err
	}
	c.Created(bucket)
	return nil
}

func (*imageCtrl) updateBucket(c *elton.Context) error {
	id, err := getIDFromParams(c)
	if err != nil {
		return err
	}
	params := bucketUpdateParams{}
	err = validateBody(c, &params)
	if err != nil {
		return err
	}
	us := getUserSession(c)
	result, err := params.updateOneID(c.Context(), id, us.MustGetInfo().Account)
	if err != nil {
		return err
	}
	c.Body = result
	return nil
}

func (*imageCtrl) listBucket(c *elton.Context) error {
	params := bucketListParams{}
	err := validateQuery(c, &params)
	if err != nil {
		return err
	}
	count := -1
	if params.ShouldCount() {
		count, err = params.count(c.Context())
		if err != nil {
			return err
		}
	}

	buckets, err := params.queryAll(c.Context())
	if err != nil {
		return err
	}
	c.Body = &bucketListResp{
		Count:   count,
		Buckets: buckets,
	}
	return nil
}

func (*imageCtrl) addImage(c *elton.Context) error {
	params := imageAddParams{
		Bucket: c.Request.FormValue("bucket"),
		Name:   c.Request.FormValue("name"),
		Tags:   c.Request.FormValue("tags"),
	}
	err := validate.Struct(&params)
	if err != nil {
		return err
	}
	if params.Name == "" {
		params.Name = util.GenXID()
	}

	account := getUserSession(c).MustGetInfo().Account
	err = validateBucketForUser(c.Context(), params.Bucket, account)
	if err != nil {
		return err
	}
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		return err
	}
	defer file.Close()
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	params.creator = account
	params.data = buf
	result, err := params.save(c.Context())
	if err != nil {
		return err
	}
	// 图片数据不返回
	result.Data = nil
	c.Created(result)
	return nil
}

func (*imageCtrl) listImage(c *elton.Context) error {
	params := imageListParams{}
	err := validateQuery(c, &params)
	if err != nil {
		return err
	}

	count := -1
	if params.ShouldCount() {
		count, err = params.count(c.Context())
		if err != nil {
			return err
		}
	}
	images, err := params.queryAll(c.Context())
	if err != nil {
		return err
	}

	c.Body = &imageListResp{
		Count:  count,
		Images: images,
	}
	return nil
}

func (*imageCtrl) getImageThumbnail(c *elton.Context) error {
	params := imageGetThumbnailParams{}
	err := validate.Query(&params, util.MergeMapString(c.Params.ToMap(), c.Query()))
	if err != nil {
		return err
	}
	jobs := []pipeline.ImageJob{
		pipeline.NewGetEntImage(params.Bucket, params.Name),
		pipeline.NewFitResizeImage(params.ThumbnailSize, params.ThumbnailSize),
	}
	img, err := pipeline.Do(c.Context(), nil, jobs...)
	if err != nil {
		return err
	}
	c.CacheMaxAge(time.Minute)
	c.SetContentTypeByExt("." + img.Type)
	c.BodyBuffer = bytes.NewBuffer(img.Data)

	return nil
}

func (*imageCtrl) pipeline(c *elton.Context) error {
	rawQuery := c.Request.URL.RawQuery
	if len(rawQuery) == 0 {
		return hes.New("pipeline can not be empty")
	}
	tasks := strings.Split(rawQuery, "|")
	jobs, err := pipeline.Parse(tasks, c.Request.Header)
	if err != nil {
		return err
	}
	img, err := pipeline.Do(c.Context(), nil, jobs...)
	if err != nil {
		return err
	}
	log.Info(c.Context()).
		Strs("tasks", tasks).
		Int("originalSize", img.OriginalSize).
		Int("size", img.Size).
		Int("percent", 100*img.Size/img.OriginalSize).
		Msg("")
	c.SetContentTypeByExt("." + img.Type)
	c.BodyBuffer = bytes.NewBuffer(img.Data)
	return nil
}
