package controller

import (
	"bytes"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"

	lru "github.com/hashicorp/golang-lru"
	"github.com/kataras/iris"
	"github.com/spf13/viper"
	"github.com/vicanso/tiny-site/global"
	"github.com/vicanso/tiny-site/middleware"
	"github.com/vicanso/tiny-site/model"
	"github.com/vicanso/tiny-site/router"
	"github.com/vicanso/tiny-site/service"
	"github.com/vicanso/tiny-site/util"
)

var (
	tmpFileCache     *lru.Cache
	errFileIsExpired = &util.HTTPError{
		StatusCode: http.StatusBadRequest,
		Category:   util.ErrCategoryValidate,
		Code:       util.ErrCodeFile,
		Message:    "the file is expired, please upload again",
	}
	errFileNotFound = &util.HTTPError{
		StatusCode: http.StatusBadRequest,
		Category:   util.ErrCategoryValidate,
		Code:       util.ErrCodeFile,
		Message:    "the file is not found",
	}
	errImageOptimOverLimit = &util.HTTPError{
		StatusCode: http.StatusBadRequest,
		Category:   util.ErrCategoryValidate,
		Code:       util.ErrCodeFile,
		Message:    "image optim over limit",
	}
	maxImageWidth   = 2048
	maxImageHeight  = 2048
	maxImageQuality = 100
	maxCacheSize    = 1024
)

type (
	fileCtrl           struct{}
	uploadInfoResponse struct {
		ID       string `json:"id,omitempty"`
		FileType string `json:"fileType,omitempty"`
	}
	saveFileParams struct {
		ID       string `json:"id,omitempty" valid:"runelength(26|26)"`
		Category string `json:"category,omitempty" valid:"runelength(2|20)"`
		FileType string `json:"fileType,omitempty" valid:"in(jpeg|png)"`
	}
	getFileParams struct {
		Type    string `json:"type,omitempty" valid:"in(jpeg|png|webp|guetzli),optional"`
		Width   string `json:"width,omitempty" valid:"int,range(0|2048),optional"`
		Height  string `json:"height,omitempty" valid:"int,range(0|2048),optional"`
		Quality string `json:"quality,omitempty" valid:"int,range(0|90),optional"`
	}
)

func init() {
	v := viper.GetInt("upload.cacheSize")
	if v != 0 {
		maxCacheSize = v
	}
	c, err := global.NewLRU(maxCacheSize)
	if err != nil {
		panic(err)
	}
	tmpFileCache = c

	v = viper.GetInt("tiny.maxWidth")
	if v != 0 {
		maxImageWidth = v
	}
	v = viper.GetInt("tiny.maxHeight")
	if v != 0 {
		maxImageHeight = v
	}
	v = viper.GetInt("tiny.maxQuality")
	if v != 0 {
		maxImageQuality = v
	}

	// TODO 添加登录校验
	files := router.NewGroup("/files")
	images := router.NewGroup("/images")

	ctrl := fileCtrl{}
	files.Add("POST", "/v1/upload", ctrl.upload)
	files.Add("POST", "/v1/save", ctrl.save)

	images.Add("GET", "/v1/:file", middleware.IsNilQuery, ctrl.get)

}

func getOptimOptions(params []string) (opts *service.OptimOptions, err error) {
	quality, err := strconv.Atoi(params[1])
	if err != nil {
		return
	}
	// file-quality-width-height.ext
	width, err := strconv.Atoi(params[2])
	if err != nil {
		return
	}
	height, err := strconv.Atoi(params[3])
	if err != nil {
		return
	}
	opts = &service.OptimOptions{
		Quality: quality,
		Width:   width,
		Height:  height,
	}
	return
}

// upload 上传文件
func (c *fileCtrl) upload(ctx iris.Context) {
	req := ctx.Request()
	file, header, err := req.FormFile("file")
	if err != nil {
		resErr(ctx, err)
		return
	}
	defer file.Close()
	fileType := filepath.Ext(header.Filename)
	if fileType != "" {
		fileType = fileType[1:]
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		resErr(ctx, err)
		return
	}
	id := util.GenUlid()
	tmpFileCache.Add(id, buf.Bytes())
	info := &uploadInfoResponse{
		ID:       id,
		FileType: fileType,
	}
	resCreated(ctx, info)
}

// save 保存文件
func (c *fileCtrl) save(ctx iris.Context) {
	params := &saveFileParams{}
	err := validate(params, getRequestBody(ctx))
	if err != nil {
		resErr(ctx, err)
		return
	}
	data, _ := tmpFileCache.Get(params.ID)
	if data == nil {
		resErr(ctx, errFileIsExpired)
		return
	}
	f := model.File{
		File: params.ID,
		Type: params.FileType,
		Data: data.([]byte),
	}
	err = f.Save()
	if err != nil {
		resErr(ctx, err)
		return
	}
	tmpFileCache.Remove(params.ID)
	resCreated(ctx, nil)
}

// get get file
func (c *fileCtrl) get(ctx iris.Context) {
	fileParam := ctx.Params().Get("file")
	ext := filepath.Ext(fileParam)
	params := strings.Split(fileParam[0:len(fileParam)-len(ext)], "-")
	if len(params) != 4 || ext == "" {
		resErr(ctx, &util.HTTPError{
			StatusCode: http.StatusBadRequest,
			Category:   util.ErrCategoryValidate,
			Code:       util.ErrCodeFile,
			Message:    "file name is wrong, it should be file-quality-width-height.ext",
		})
		return
	}
	file := params[0]

	f := model.File{
		File: file,
	}
	err := f.First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = errFileNotFound
		}
		resErr(ctx, err)
		return
	}
	opts, err := getOptimOptions(params)
	if err != nil {
		resErr(ctx, err)
		return
	}
	opts.Type = ext[1:]
	opts.Data = f.Data
	if opts.Width > maxImageWidth ||
		opts.Height > maxImageHeight ||
		opts.Quality > maxImageQuality {
		resErr(ctx, errImageOptimOverLimit)
		return
	}

	buf, err := service.Optim(opts)
	if err != nil {
		resErr(ctx, err)
		return
	}
	// 30 day
	setCache(ctx, "720h")
	// convert ext
	if opts.Type == "guetzli" {
		ext = ".jpeg"
	}
	ctx.ContentType(mime.TypeByExtension(ext))
	res(ctx, buf)
}
