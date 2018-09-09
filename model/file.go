package model

import (
	"net/http"

	"github.com/vicanso/tiny-site/util"
)

type (
	// File file model
	File struct {
		BaseModel
		File     string `json:"file,omitempty" gorm:"type:varchar(26);not null;unique_index:idx_files_file"`
		Category string `json:"category,omitempty"`
		Type     string `json:"type,omitempty"`
		Data     []byte `json:"data,omitempty"`
	}
)

var (
	errFileFieldIsNil = &util.HTTPError{
		StatusCode: http.StatusBadRequest,
		Category:   util.ErrCategoryLogic,
		Code:       util.ErrCodeFile,
		Message:    "file, fileType and data cant not be nil",
	}
)

// Save save file
func (f *File) Save() (err error) {
	if f.File == "" || f.Type == "" || len(f.Data) == 0 {
		return errFileFieldIsNil
	}
	client := GetClient()
	err = client.Create(f).Error
	return
}

// First get file
func (f *File) First() (err error) {
	client := GetClient()
	err = client.Where(f).First(f).Error
	return
}
