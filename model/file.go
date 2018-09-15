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
		MaxAge   string `json:"maxAge,omitempty" gorm:"type:varchar(10)"`
		Category string `json:"category,omitempty" gorm:"type:varchar(20)"`
		Type     string `json:"type,omitempty" gorm:"type:varchar(10)"`
		Size     int    `json:"size,omitempty"`
		Data     []byte `json:"data,omitempty"`
		Creator  string `json:"creator,omitempty" gorm:"type:varchar(20);not null"`
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

// List list the file
func (f *File) List(fields, order string, skip, limit int) (files []*File, err error) {
	client := GetClient()
	files = make([]*File, 0)
	c := client.Where(f)
	if fields != "" {
		c = c.Select(convertFields(fields))
	}
	if order != "" {
		c = c.Order(convertOrder(order))
	}
	if skip != 0 {
		c = c.Offset(skip)
	}
	if limit != 0 {
		c = c.Limit(limit)
	}
	err = c.Find(&files).Error
	return
}

// Count count the file
func (f *File) Count() (count int, err error) {
	client := GetClient()
	err = client.Model(&File{}).Where(f).Count(&count).Error
	return
}

// GetCategories get all category
func (f *File) GetCategories() (categories []string, err error) {
	client := GetClient()
	rows, err := client.Raw("SELECT DISTINCT category FROM files").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		v := ""
		rows.Scan(&v)
		if v != "" {
			categories = append(categories, v)
		}
	}
	return
}
