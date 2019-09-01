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

import "time"

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

		User      string `json:"user,omitempty" gorm:"type:varchar(20);not null;unique_index:idx_file_zone_authorities_account"`
		Authority int    `json:"authority,omitempty"`
		Zone      int    `json:"zone,omitempty" gorm:"not null"`
	}
)

const (
	// AuthorityNone none authority
	AuthorityNone = iota
	// AuthorityRead read authority
	AuthorityRead
	// AuthorityReadWrite read/write authority
	AuthorityReadWrite
)

func init() {
	pgGetClient().AutoMigrate(&File{}).
		AutoMigrate(&FileZone{}).
		AutoMigrate(&FileZoneAuthority{})
}
