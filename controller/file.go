package controller

import (
	"bytes"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"

	lru "github.com/hashicorp/golang-lru"
	"github.com/kataras/iris"
	"github.com/spf13/viper"
	"github.com/vicanso/tiny-site/cs"
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
		Category string `json:"category,omitempty" valid:"runelength(2|20)"`
		FileType string `json:"fileType,omitempty" valid:"in(jpeg|png)"`
		MaxAge   string `json:"maxAge,omitempty" valid:"matches(^[0-9]+[smh]$)"`
	}
	listFileParams struct {
		Category string `json:"category,omitempty" valid:"runelength(2|20)"`
		Fields   string `json:"fields" valid:"runelength(2|100)"`
		Order    string `json:"order" valid:"optional"`
		Skip     string `json:"skip" valid:"int,optional"`
		Limit    string `json:"limit" valid:"in(1|10|30),optional"`
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

	files := router.NewGroup("/files")
	images := router.NewGroup("/images")

	ctrl := fileCtrl{}
	// TODO 添加登录校验
	files.Add("POST", "/v1/upload", ctrl.upload)
	// TODO 添加登录校验
	files.Add(
		"POST",
		"/v1/:id",
		router.SessionHandler,
		middleware.IsLogined,
		ctrl.save,
	)

	files.Add("GET", "/v1/categories", ctrl.getCategories)
	files.Add("GET", "/v1", ctrl.list)

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
	id := ctx.Params().Get("id")
	params := &saveFileParams{}
	err := validate(params, getRequestBody(ctx))
	if err != nil {
		resErr(ctx, err)
		return
	}
	data, _ := tmpFileCache.Get(id)
	if data == nil {
		resErr(ctx, errFileIsExpired)
		return
	}
	buf := data.([]byte)
	sess := getSession(ctx)
	f := model.File{
		File:     id,
		Type:     params.FileType,
		Data:     buf,
		Category: params.Category,
		Size:     len(buf),
		MaxAge:   params.MaxAge,
		Creator:  sess.GetString(cs.SessionAccountField),
	}
	err = f.Save()
	if err != nil {
		resErr(ctx, err)
		return
	}
	tmpFileCache.Remove(id)
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
	if f.MaxAge != "" {
		setCache(ctx, f.MaxAge)
	}

	// convert ext
	if opts.Type == "guetzli" {
		ext = ".jpeg"
	}
	ctx.ContentType(mime.TypeByExtension(ext))
	res(ctx, buf)
}

// getCategories get the category
func (c *fileCtrl) getCategories(ctx iris.Context) {
	f := &model.File{}
	categories, err := f.GetCategories()
	if err != nil {
		resErr(ctx, err)
		return
	}
	sort.Slice(categories, func(i, j int) bool {
		return strings.Compare(categories[i], categories[j]) < 0
	})
	m := map[string]interface{}{
		"categories": categories,
	}
	setCache(ctx, "1m")
	res(ctx, m)
}

// list list the files
func (c *fileCtrl) list(ctx iris.Context) {
	params := &listFileParams{}
	err := validate(params, getRequestQuery(ctx))
	if err != nil {
		resErr(ctx, err)
		return
	}
	f := &model.File{
		Category: params.Category,
	}
	files, err := f.List(params.Fields, params.Order)
	if err != nil {
		resErr(ctx, err)
		return
	}
	m := map[string]interface{}{
		"files": files,
	}
	res(ctx, m)
}
