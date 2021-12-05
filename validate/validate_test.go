package validate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type subData struct {
	SubTitle string `json:"subTitle"`
}

type Data struct {
	Title   string   `json:"title"`
	SubData *subData `json:"subData"`
}

type MergeData struct {
	subData
	Date  time.Time `json:"date"`
	Title string    `json:"title"`
	Name  string    `json:"name" default:"test"`
}

func TestValidateQuery(t *testing.T) {
	assert := assert.New(t)
	md := MergeData{}
	data := map[string]string{
		"date":     "2021-10-23T01:33:53.344Z",
		"subTitle": "s",
		"title":    "t",
	}
	err := Query(&md, data)
	assert.Nil(err)
	assert.Equal("2021-10-23 01:33:53.344 +0000 UTC", md.Date.String())
	assert.Equal("s", md.SubTitle)
	assert.Equal("s", md.SubTitle)
	assert.Equal("t", md.Title)
	assert.Equal("test", md.Name)
}

func TestValidate(t *testing.T) {
	assert := assert.New(t)
	md := MergeData{}
	data := map[string]string{
		"subTitle": "s",
		"title":    "t",
		"name":     "123",
	}
	err := Do(&md, data)
	assert.Nil(err)
	assert.Equal("s", md.SubTitle)
	assert.Equal("s", md.SubTitle)
	assert.Equal("t", md.Title)
	assert.Equal("123", md.Name)

	d := Data{}
	err = Do(&d, []byte(`{
		"title": "t",
		"subData": {
			"subTitle": "s"
		}
	}`))
	assert.Nil(err)
	assert.Equal("s", d.SubData.SubTitle)
	assert.Equal("t", d.Title)
}
