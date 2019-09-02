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

	"github.com/jinzhu/gorm"
)

type (
	// File file struct
	File struct {
		ID        uint       `gorm:"primary_key" json:"id,omitempty"`
		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt time.Time  `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `sql:"index" json:"deletedAt,omitempty"`

		Name    string `json:"name,omitempty" gorm:"type:varchar(26);not null;unique_index:idx_files_name"`
		MaxAge  string `json:"maxAge,omitempty" gorm:"type:varchar(10)"`
		Zone    int    `json:"zone,omitempty" gorm:"not null"`
		Type    string `json:"type,omitempty" gorm:"type:varchar(10)"`
		Size    int    `json:"size,omitempty"`
		Width   int    `json:"width,omitempty"`
		Height  int    `json:"height,omitempty"`
		Data    []byte `json:"data,omitempty"`
		Creator string `json:"creator,omitempty" gorm:"type:varchar(20);not null"`
	}
	// FileZone file zone
	FileZone struct {
		ID        uint       `gorm:"primary_key" json:"id,omitempty"`
		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt time.Time  `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `sql:"index" json:"deletedAt,omitempty"`

		Name  string `json:"name,omitempty" gorm:"type:varchar(26);not null;unique_index:idx_file_zones_name"`
		Owner string `json:"owner,omitempty" gorm:"type:varchar(20);not null"`
	}
	// FileZoneAuthority file zone authority
	FileZoneAuthority struct {
		ID        uint       `gorm:"primary_key" json:"id,omitempty"`
		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt time.Time  `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `sql:"index" json:"deletedAt,omitempty"`

		User      string `json:"user,omitempty" gorm:"type:varchar(20);not null;unique_index:idx_file_zone_authorities_user_zone"`
		Authority int    `json:"authority,omitempty" `
		Zone      int    `json:"zone,omitempty" gorm:"not null;unique_index:idx_file_zone_authorities_user_zone"`
	}
	// FileZoneQueryParams file zone query params
	FileZoneQueryParams struct {
	}
	// FileSrv file service
	FileSrv struct{}
)

const (
	// AuthorityNone none authority
	AuthorityNone = iota
	// AuthorityRead read authority
	AuthorityRead
	// AuthorityReadWrite read/write authority
	AuthorityReadWrite
)

var ()

func init() {
	pgGetClient().AutoMigrate(&File{}).
		AutoMigrate(&FileZone{}).
		AutoMigrate(&FileZoneAuthority{})
}

// Add add file
func (srv *FileSrv) Add(f *File) (err error) {
	err = pgCreate(f)
	return
}

// AddZone add file zone
func (srv *FileSrv) AddZone(fz *FileZone) (err error) {
	err = pgCreate(fz)
	return
}

// ListZone list file zone
func (srv *FileSrv) ListZone() (result []*FileZone, err error) {
	result = make([]*FileZone, 0)
	db := pgGetClient()
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

// AddZoneAuthority add file zone authority
func (srv *FileSrv) AddZoneAuthority(fza *FileZoneAuthority) (err error) {
	err = pgCreate(fza)
	return
}

// UpdateZoneAuthority update file zone authority
func (srv *FileSrv) UpdateZoneAuthority(fza *FileZoneAuthority, attrs ...interface{}) (err error) {
	err = pgGetClient().Model(fza).Update(attrs...).Error
	return
}

// GetZoneAuthority get file zone authority
func (srv *FileSrv) GetZoneAuthority(conditions *FileZoneAuthority) (fza *FileZoneAuthority, err error) {
	fza = &FileZoneAuthority{}
	err = pgGetClient().First(fza, conditions).Error
	return
}

// ZoneWritable check zone wriable
func (srv *FileSrv) ZoneWritable(user string, zone int) (writable bool, err error) {
	fza, err := srv.GetZoneAuthority(&FileZoneAuthority{
		User: user,
		Zone: zone,
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
		}
		return
	}
	if fza.Authority == AuthorityReadWrite {
		writable = true
	}
	return
}

// DeleteZoneAuthorityByID delete file zone authority
func (srv *FileSrv) DeleteZoneAuthorityByID(id uint) (err error) {
	err = pgGetClient().Unscoped().Delete(&FileZoneAuthority{
		ID: id,
	}).Error
	return
}
