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
	"time"
)

type (
	// File file struct
	File struct {
		ID        uint       `gorm:"primary_key" json:"id,omitempty"`
		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt time.Time  `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `sql:"index" json:"deletedAt,omitempty"`

		Name        string `json:"name,omitempty" gorm:"type:varchar(26);not null;unique_index:idx_files_name"`
		Description string `json:"description,omitempty"`
		MaxAge      string `json:"maxAge,omitempty" gorm:"type:varchar(10)"`
		Zone        int    `json:"zone,omitempty" gorm:"not null"`
		Type        string `json:"type,omitempty" gorm:"type:varchar(10)"`
		Size        int    `json:"size,omitempty"`
		Width       int    `json:"width,omitempty"`
		Height      int    `json:"height,omitempty"`
		Data        []byte `json:"data,omitempty"`
		Thumbnail   []byte `json:"thumbnail,omitempty"`
		Creator     string `json:"creator,omitempty" gorm:"type:varchar(20);not null"`
	}
	// FileZone file zone
	FileZone struct {
		ID        uint       `gorm:"primary_key" json:"id,omitempty"`
		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt time.Time  `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `sql:"index" json:"deletedAt,omitempty"`

		Name        string `json:"name,omitempty" gorm:"type:varchar(26);not null;unique_index:idx_file_zones_name"`
		Description string `json:"description,omitempty"`
		Owner       string `json:"owner,omitempty" gorm:"type:varchar(20);not null"`
	}
	// FileQueryParams file query params
	FileQueryParams struct {
		Keyword string
		Limit   int
		Offset  int
		Zone    int
		Fields  string
		Sort    string
	}
	// FileSrv file service
	FileSrv struct{}
)

func init() {
	pgGetClient().AutoMigrate(&File{}).
		AutoMigrate(&FileZone{})
}

// List list file
func (srv *FileSrv) List(params FileQueryParams) (result []*File, err error) {
	result = make([]*File, 0)
	db := pgGetClient()
	if params.Limit > 0 {
		db = db.Limit(params.Limit)
	}
	if params.Offset > 0 {
		db = db.Offset(params.Offset)
	}
	if params.Fields != "" {
		db = db.Select(pgFormatSelect(params.Fields))
	}
	if params.Sort != "" {
		db = db.Order(pgFormatOrder(params.Sort))
	}
	db = db.Where("zone = (?)", params.Zone)
	if params.Keyword != "" {
		db = db.Where("name LIKE ?", "%"+params.Keyword+"%").
			Or("description LIKE ?", "%"+params.Keyword+"%")
	}
	err = db.Find(&result).Error
	return
}

// Count count file
func (srv *FileSrv) Count(params FileQueryParams) (count int, err error) {
	db := pgGetClient().Model(&File{})
	err = db.Where("zone = (?)", params.Zone).Count(&count).Error
	return
}

// Add add file
func (srv *FileSrv) Add(f *File) (err error) {
	err = pgCreate(f)
	return
}

// GetByName get file by name
func (srv *FileSrv) GetByName(name string) (f *File, err error) {
	f = &File{}
	err = pgGetClient().First(f, File{
		Name: name,
	}).Error
	return
}

// FindByID get file by id
func (srv *FileSrv) FindByID(id uint, args ...string) (f *File, err error) {
	fields := ""
	if len(args) > 0 {
		fields = args[0]
	}
	f = &File{}
	err = pgGetClient().Select(pgFormatSelect(fields)).First(f, File{
		ID: id,
	}).Error
	return
}

// UpdateByID update by id
func (srv *FileSrv) UpdateByID(id uint, f *File) (err error) {
	err = pgGetClient().Model(&File{
		ID: id,
	}).Update(f).Error
	return
}

// AddZone add file zone
func (srv *FileSrv) AddZone(fz *FileZone) (err error) {
	err = pgCreate(fz)
	return
}

// ListZone list file zone
func (srv *FileSrv) ListZone(conditions *FileZone) (result []*FileZone, err error) {
	result = make([]*FileZone, 0)
	db := pgGetClient()
	if conditions != nil {
		db = db.Where(conditions)
	}
	err = db.Find(&result).Error
	return
}

// GetZone get file zone
func (srv *FileSrv) GetZone(conditions *FileZone) (fz *FileZone, err error) {
	fz = &FileZone{}
	err = pgGetClient().First(fz, conditions).Error
	return
}

// UpdateZone update fie zone
func (srv *FileSrv) UpdateZone(fz *FileZone, attrs ...interface{}) (err error) {
	err = pgGetClient().Model(fz).Update(attrs...).Error
	return
}
