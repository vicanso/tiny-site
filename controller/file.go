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
		Name   string `json:"name,omitempty" valid:"xFileName"`
		MaxAge string `json:"maxAge,omitempty" valid:"xDuration"`
		Zone   int    `json:"zone,omitempty" valid:"xFileZone"`
		Type   string `json:"type,omitempty" valid:"xFileType"`
		Width  int    `json:"width,omitempty" valid:"-"`
		Height int    `json:"height,omitempty" valid:"-"`
		Data   string `json:"data,omitempty" valid:"-"`
	}
	createFileZoneParams struct {
		Name  string `json:"name,omitempty" valid:"xFileZoneName"`
		Owner string `json:"owner,omitempty" valid:"xUserAccount"`
	}
	transferFileZoneParams struct {
		TransferTo string `json:"transferTo,omitempty" valid:"xUserAccount"`
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
	errNotAllowToAdjustFileZone = &hes.Error{
		Message:    "not allow to adjust file zone",
		StatusCode: http.StatusForbidden,
	}
	errNoWriteAuthority = &hes.Error{
		Message:    "not allow to add file to zone",
		StatusCode: http.StatusForbidden,
	}
	errFileZoneIDInvalid = hes.New("file zone id is invalid")
)

const (
	fileZoneIDKey          = "fileZoneID"
	fileZoneAuthorityIDKey = "fileZoneAuthorityID"
)

func init() {
	ctrl := fileCtrl{}
	g := router.NewGroup("/files")

	// 创建文件
	g.POST("/v1/upload/save", shouldLogined, ctrl.create)

	g.POST("/v1/upload", ctrl.upload)

	g.GET("/v1/zones", shouldLogined, ctrl.listZone)

	// 创建file zone，只允许admin权限用户创建
	g.POST("/v1/zones", shouldBeAdmin, ctrl.createZone)

	// file zone拥有者转移
	g.POST(
		"/v1/zones/:fileZoneID/transfer",
		shouldLogined,
		shouldAdminOrFileZoneOwner,
		newTracker(cs.ActionFileZoneTransfer),
		ctrl.transferZone,
	)

	// 添加file zone的权限
	g.POST(
		"/v1/zones/:fileZoneID/authorities",
		shouldLogined,
		shouldAdminOrFileZoneOwner,
		newTracker(cs.ActionFileZoneAuthorityAdd),
		ctrl.createZoneAuthority,
	)
	// 更新file zone权限
	g.PATCH(
		"/v1/zones/:fileZoneID/authorities/:fileZoneAuthorityID",
		shouldLogined,
		shouldAdminOrFileZoneOwner,
		shouldAllowToUpdateFIleZoneAuthority,
		newTracker(cs.ActionFileZoneAuthorityUpdate),
		ctrl.updateZoneAuthority,
	)
	// 删除file zone权限
	g.DELETE(
		"/v1/zones/:fileZoneID/authorities/:fileZoneAuthorityID",
		shouldLogined,
		shouldAdminOrFileZoneOwner,
		shouldAllowToUpdateFIleZoneAuthority,
		newTracker(cs.ActionFileZoneAuthorityDelete),
		ctrl.deleteZoneAuthority,
	)
}

func shouldAdminOrFileZoneOwner(c *elton.Context) (err error) {
	id, _ := strconv.Atoi(c.Param(fileZoneIDKey))
	us := service.NewUserSession(c)
	roles := us.GetRoles()
	// 如果当前用户不是管理员，需判断是否该空间的owner
	if !util.UserRoleIsValid(adminUserRoles, roles) {
		fz, err := fileSrv.GetZone(&service.FileZone{
			ID: uint(id),
		})
		if err != nil {
			return err
		}
		if us.GetAccount() != fz.Owner {
			return errNotAllowToAdjustFileZone
		}
	}
	return c.Next()
}

func shouldAllowToUpdateFIleZoneAuthority(c *elton.Context) (err error) {
	fileZoneID, _ := strconv.Atoi(c.Param(fileZoneIDKey))
	fileZoneAuthorityID, _ := strconv.Atoi(c.Param(fileZoneAuthorityIDKey))
	fza, err := fileSrv.GetZoneAuthority(&service.FileZoneAuthority{
		ID: uint(fileZoneAuthorityID),
	})
	if err != nil {
		return
	}
	if fza.Zone != fileZoneID {
		err = errFileZoneIDInvalid
		return
	}
	return c.Next()
}

func (ctrl fileCtrl) create(c *elton.Context) (err error) {
	params := &createFileParams{}
	err = validate.Do(params, c.RequestBody)
	if err != nil {
		return
	}
	// 判断用户是否有该zone的写权限
	us := service.NewUserSession(c)
	account := us.GetAccount()
	writable, err := fileSrv.ZoneWritable(account, params.Zone)
	if err != nil {
		return
	}
	if !writable {
		err = errNoWriteAuthority
		return
	}
	buf, err := base64.StdEncoding.DecodeString(params.Data)
	if err != nil {
		return
	}
	if len(buf) == 0 {
		err = hes.New("data can't be nil")
		return
	}
	f := &service.File{
		Name:    params.Name,
		MaxAge:  params.MaxAge,
		Zone:    params.Zone,
		Type:    params.Type,
		Size:    len(buf),
		Width:   params.Width,
		Height:  params.Height,
		Data:    buf,
		Creator: account,
	}
	err = fileSrv.Add(f)
	if err != nil {
		return
	}
	// 响应数据时把data清除
	f.Data = nil
	c.Created(f)
	return
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

func (ctrl fileCtrl) createZone(c *elton.Context) (err error) {
	params := &createFileZoneParams{}
	err = validate.Do(params, c.RequestBody)
	if err != nil {
		return
	}
	fz := &service.FileZone{
		Name:  params.Name,
		Owner: params.Owner,
	}
	err = fileSrv.AddZone(fz)
	if err != nil {
		return
	}
	c.Created(fz)
	return
}

func (ctrl fileCtrl) listZone(c *elton.Context) (err error) {
	result, err := fileSrv.ListZone()
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

func (ctrl fileCtrl) transferZone(c *elton.Context) (err error) {
	params := &transferFileZoneParams{}
	err = validate.Do(params, c.RequestBody)
	if err != nil {
		return
	}
	id, _ := strconv.Atoi(c.Param(fileZoneIDKey))
	err = fileSrv.UpdateZone(&service.FileZone{
		ID: uint(id),
	}, service.FileZone{
		Owner: params.TransferTo,
	})
	if err != nil {
		return
	}
	c.NoContent()
	return
}

func (ctrl fileCtrl) createZoneAuthority(c *elton.Context) (err error) {
	params := &createFileZoneAuthorityParams{}
	err = validate.Do(params, c.RequestBody)
	if err != nil {
		return
	}
	fileZoneID, _ := strconv.Atoi(c.Param(fileZoneIDKey))
	fza := &service.FileZoneAuthority{
		User:      params.User,
		Authority: params.Authority,
		Zone:      fileZoneID,
	}
	err = fileSrv.AddZoneAuthority(fza)
	if err != nil {
		return
	}
	c.Created(fza)
	return
}

func (ctrl fileCtrl) updateZoneAuthority(c *elton.Context) (err error) {
	params := &updateFileZoneAuthorityParams{}
	err = validate.Do(params, c.RequestBody)
	if err != nil {
		return
	}
	fileZoneAuthorityID, _ := strconv.Atoi(c.Param(fileZoneAuthorityIDKey))
	err = fileSrv.UpdateZoneAuthority(&service.FileZoneAuthority{
		ID: uint(fileZoneAuthorityID),
	}, &service.FileZoneAuthority{
		Authority: params.Authority,
	})
	if err != nil {
		return
	}
	c.NoContent()
	return
}

func (ctrl fileCtrl) deleteZoneAuthority(c *elton.Context) (err error) {
	fileZoneAuthorityID, _ := strconv.Atoi(c.Param(fileZoneAuthorityIDKey))
	err = fileSrv.DeleteZoneAuthorityByID(uint(fileZoneAuthorityID))
	if err != nil {
		return
	}
	c.NoContent()
	return
}
