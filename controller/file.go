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
	"encoding/base64"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/vicanso/elton"
	"github.com/vicanso/hes"
	"github.com/vicanso/tiny-site/cs"
	"github.com/vicanso/tiny-site/router"
	"github.com/vicanso/tiny-site/service"
	"github.com/vicanso/tiny-site/util"
	"github.com/vicanso/tiny-site/validate"
)

type (
	fileCtrl struct{}
	fileInfo struct {
		Name   string `json:"name,omitempty"`
		Data   []byte `json:"data,omitempty"`
		Type   string `json:"type,omitempty"`
		Size   int    `json:"size,omitempty"`
		Width  int    `json:"width,omitempty"`
		Height int    `json:"height,omitempty"`
	}
	createFileParams struct {
		Name        string `json:"name,omitempty" valid:"xFileName"`
		Description string `json:"description,omitempty" valid:"xFileDesc"`
		MaxAge      string `json:"maxAge,omitempty" valid:"xDuration"`
		Zone        int    `json:"zone,omitempty" valid:"xFileZone"`
		Type        string `json:"type,omitempty" valid:"xFileType"`
		Width       int    `json:"width,omitempty" valid:"-"`
		Height      int    `json:"height,omitempty" valid:"-"`
		Data        string `json:"data,omitempty" valid:"-"`
	}
	updateFileParams struct {
		Description string `json:"description,omitempty" valid:"xFileDesc,optional"`
		MaxAge      string `json:"maxAge,omitempty" valid:"xDuration,optional"`
		Type        string `json:"type,omitempty" valid:"xFileType,optional"`
		Width       int    `json:"width,omitempty" valid:"-"`
		Height      int    `json:"height,omitempty" valid:"-"`
		Data        string `json:"data,omitempty" valid:"-"`
	}
	listFileParams struct {
		Zone   string `json:"zone,omitempty" valid:"xFileZone"`
		Limit  string `json:"limit,omitempty" valid:"xLimit"`
		Offset string `json:"offset,omitempty" valid:"xOffset"`
		Fields string `json:"fields,omitempty" valid:"xFields"`
		Sort   string `json:"sort,omitempty" valid:"xSort,optional"`
	}
	createFileZoneParams struct {
		Name        string `json:"name,omitempty" valid:"xFileZoneName"`
		Description string `json:"description,omitempty" valid:"xFileZoneDesc"`
		Owner       string `json:"owner,omitempty" valid:"xUserAccount"`
	}
	updateFileZoneParams struct {
		Name        string `json:"name,omitempty" valid:"xFileZoneName,optional"`
		Description string `json:"description,omitempty" valid:"xFileZoneDesc,optional"`
		Owner       string `json:"owner,omitempty" valid:"xUserAccount,optional"`
	}
	createFileZoneAuthorityParams struct {
		User      string `json:"user,omitempty" valid:"xUserAccount"`
		Authority int    `json:"authority,omitempty" valid:"xFileZoneAuthority"`
	}
	updateFileZoneAuthorityParams struct {
		Authority int `json:"authority,omitempty" valid:"xFileZoneAuthority"`
	}
)

var (
	errNotAllowToUpdateFileZone = &hes.Error{
		Message:    "not allow to update file zone",
		StatusCode: http.StatusForbidden,
	}
	errNoWriteAuthority = &hes.Error{
		Message:    "not allow to add file to zone",
		StatusCode: http.StatusForbidden,
	}
	errNoReadAuthority = &hes.Error{
		Message:    "not allow to read file from zone",
		StatusCode: http.StatusForbidden,
	}
	errNotAllowToUpdateFile = &hes.Error{
		Message:    "not allow to update file",
		StatusCode: http.StatusForbidden,
	}
	errFileDataIsNil     = hes.New("data can't be nil")
	errFileZoneIDInvalid = hes.New("file zone id is invalid")
)

const (
	fileIDKey              = "fileID"
	fileZoneIDKey          = "fileZoneID"
	fileZoneAuthorityIDKey = "fileZoneAuthorityID"
)

const (
	thumbnailWidth   = 60
	thumbnailQuality = 70
)

func init() {
	ctrl := fileCtrl{}
	g := router.NewGroup("/files")

	// 获取文件列表
	g.GET("/v1", shouldLogined, ctrl.list)
	// 获取文件详情
	g.GET("/v1/detail/:fileID", shouldLogined, ctrl.detail)
	// 创建文件
	g.POST("/v1/upload/save", shouldLogined, ctrl.create)
	// 上传文件
	g.POST("/v1/upload", shouldLogined, ctrl.upload)
	// 更新文件
	g.PATCH("/v1/upload/:fileID", shouldLogined, ctrl.updateUpload)

	// 获取文件空间列表
	g.GET("/v1/zones", shouldLogined, ctrl.listZone)
	// 获取我的文件空间列表
	g.GET("/v1/zones/mine", shouldLogined, ctrl.listMyZone)

	// 创建file zone，只允许admin权限用户创建
	g.POST("/v1/zones", shouldBeAdmin, ctrl.createZone)

	// file zone更新
	g.PATCH(
		"/v1/zones/:fileZoneID",
		shouldLogined,
		shouldAdminOrFileZoneOwner,
		newTracker(cs.ActionFileZoneUpdate),
		ctrl.updateZone,
	)

}

func shouldAdminOrFileZoneOwner(c *elton.Context) (err error) {
	id, _ := strconv.Atoi(c.Param(fileZoneIDKey))
	us := service.NewUserSession(c)
	// 如果当前用户不是管理员，需判断是否该空间的owner
	if !us.IsAdmin() {
		fz, err := fileSrv.GetZone(&service.FileZone{
			ID: uint(id),
		})
		if err != nil {
			return err
		}
		if us.GetAccount() != fz.Owner {
			return errNotAllowToUpdateFileZone
		}
	}
	return c.Next()
}

func (ctrl fileCtrl) create(c *elton.Context) (err error) {
	params := &createFileParams{}
	err = validate.Do(params, c.RequestBody)
	if err != nil {
		return
	}
	us := service.NewUserSession(c)
	account := us.GetAccount()

	buf, err := base64.StdEncoding.DecodeString(params.Data)
	if err != nil {
		return
	}
	if len(buf) == 0 {
		err = errFileDataIsNil
		return
	}

	thumbnail, err := ctrl.createThumbnail(buf, params.Type)
	if err != nil {
		return
	}

	f := &service.File{
		Name:        params.Name,
		MaxAge:      params.MaxAge,
		Zone:        params.Zone,
		Type:        params.Type,
		Size:        len(buf),
		Width:       params.Width,
		Height:      params.Height,
		Data:        buf,
		Thumbnail:   thumbnail,
		Description: params.Description,
		Creator:     account,
	}
	err = fileSrv.Add(f)
	if err != nil {
		return
	}
	// 响应数据时把data清除，节约带宽
	f.Data = nil
	c.Created(f)
	return
}

func (ctrl fileCtrl) upload(c *elton.Context) (err error) {
	file, header, err := c.Request.FormFile("file")
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
		Name: util.GenUlid(),
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

func (ctrl fileCtrl) createThumbnail(data []byte, t string) ([]byte, error) {
	return optimSrv.Image(service.ImageOptimParams{
		Data:       data,
		Type:       t,
		SourceType: t,
		Quality:    thumbnailQuality,
		Width:      thumbnailWidth,
	})
}

func (ctrl fileCtrl) updateUpload(c *elton.Context) (err error) {
	id, err := strconv.Atoi(c.Param(fileIDKey))
	if err != nil {
		return
	}
	params := &updateFileParams{}
	err = validate.Do(params, c.RequestBody)
	if err != nil {
		return
	}
	f, err := fileSrv.FindByID(uint(id), "creator")
	if err != nil {
		return
	}
	us := service.NewUserSession(c)
	account := us.GetAccount()
	if f.Creator != account {
		err = errNotAllowToUpdateFile
		return
	}
	buf, err := base64.StdEncoding.DecodeString(params.Data)
	if err != nil {
		return
	}
	thumbnail, err := ctrl.createThumbnail(buf, params.Type)
	if err != nil {
		return
	}

	err = fileSrv.UpdateByID(uint(id), &service.File{
		Description: params.Description,
		MaxAge:      params.MaxAge,
		Type:        params.Type,
		Width:       params.Width,
		Height:      params.Height,
		Data:        buf,
		Thumbnail:   thumbnail,
		Size:        len(buf),
	})
	if err != nil {
		return
	}
	c.NoContent()
	return
}

func (ctrl fileCtrl) list(c *elton.Context) (err error) {
	params := &listFileParams{}
	err = validate.Do(params, c.Query())
	if err != nil {
		return
	}
	zone, _ := strconv.Atoi(params.Zone)
	limit, _ := strconv.Atoi(params.Limit)
	offset, _ := strconv.Atoi(params.Offset)

	queryParams := service.FileQueryParams{
		Limit:  limit,
		Zone:   zone,
		Offset: offset,
		Fields: params.Fields,
		Sort:   params.Sort,
	}
	result, err := fileSrv.List(queryParams)
	if err != nil {
		return
	}
	count := -1
	if offset == 0 {
		count, err = fileSrv.Count(queryParams)
	}
	c.Body = &struct {
		Files []*service.File `json:"files,omitempty"`
		Count int             `json:"count,omitempty"`
	}{
		result,
		count,
	}
	return
}

func (ctrl fileCtrl) detail(c *elton.Context) (err error) {
	id, _ := strconv.Atoi(c.Param(fileIDKey))
	f, err := fileSrv.FindByID(uint(id), c.QueryParam("fields"))
	if err != nil {
		return err
	}
	c.Body = f
	return
}

func (ctrl fileCtrl) createZone(c *elton.Context) (err error) {
	params := &createFileZoneParams{}
	err = validate.Do(params, c.RequestBody)
	if err != nil {
		return
	}
	fz := &service.FileZone{
		Name:        params.Name,
		Owner:       params.Owner,
		Description: params.Description,
	}
	err = fileSrv.AddZone(fz)
	if err != nil {
		return
	}
	c.Created(fz)
	return
}

func (ctrl fileCtrl) listMyZone(c *elton.Context) (err error) {
	var conditions *service.FileZone
	us := service.NewUserSession(c)
	if !us.IsAdmin() {
		conditions = &service.FileZone{
			Owner: us.GetAccount(),
		}
	}
	result, err := fileSrv.ListZone(conditions)
	if err != nil {
		return
	}
	c.Body = &struct {
		FileZones []*service.FileZone `json:"fileZones,omitempty"`
	}{
		result,
	}
	return
}

func (ctrl fileCtrl) listZone(c *elton.Context) (err error) {
	result, err := fileSrv.ListZone(nil)
	if err != nil {
		return
	}
	c.Body = &struct {
		FileZones []*service.FileZone `json:"fileZones,omitempty"`
	}{
		result,
	}
	return
}

func (ctrl fileCtrl) updateZone(c *elton.Context) (err error) {
	params := &updateFileZoneParams{}
	err = validate.Do(params, c.RequestBody)
	if err != nil {
		return
	}
	id, _ := strconv.Atoi(c.Param(fileZoneIDKey))
	err = fileSrv.UpdateZone(&service.FileZone{
		ID: uint(id),
	}, service.FileZone{
		Name:        params.Name,
		Description: params.Description,
		Owner:       params.Owner,
	})
	if err != nil {
		return
	}
	c.NoContent()
	return
}
