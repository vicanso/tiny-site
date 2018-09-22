package middleware

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"go.uber.org/zap"

	"github.com/kataras/iris"
	"github.com/vicanso/fresh"
	"github.com/vicanso/tiny-site/asset"
	"github.com/vicanso/tiny-site/cs"
	"github.com/vicanso/tiny-site/util"
)

const (
	defaultIndexFile = "index.html"
	minCompressSize  = 1024
)

var (
	textRegexp = regexp.MustCompile("text|javascript|json")
)

type (
	// StaticServeConfig static serve config
	StaticServeConfig struct {
		// Path the static file path
		Path string
		// Compression compress the file use gzip
		Compression bool
		// MaxAge max age for cache
		MaxAge string
		// SMaxAge s-maxage for cache
		SMaxAge string
		// ETag enable etag
		ETag bool
		// Header header for static file
		Header map[string]string
		// Asset static asset
		Asset *asset.Asset
	}
)

func doFresh(ctx iris.Context, etag []byte) (done bool) {
	ifNoneMatch := ctx.GetHeader(cs.HeaderIfNoneMatch)
	// not modified (check etag)
	if fresh.Check(nil, []byte(ifNoneMatch), nil, nil, etag) {
		util.RemoveHeader(ctx, cs.HeaderCacheControl)
		ctx.WriteNotModified()
		done = true
	}
	return
}

func serveFile(ctx iris.Context, filename string, c *StaticServeConfig) (err error) {
	compression := c.Compression
	f, e := os.Open(filename)
	if e != nil {
		resErr(ctx, &util.HTTPError{
			StatusCode: http.StatusNotFound,
			Category:   util.ErrCategoryLogic,
			Code:       util.ErrCodeValidate,
			Message:    e.Error(),
		})
		return
	}
	defer f.Close()
	fi, _ := f.Stat()
	if fi.IsDir() {
		file := path.Join(filename, defaultIndexFile)
		return serveFile(ctx, file, c)
	}
	size := fi.Size()
	name := fi.Name()
	contentType := mime.TypeByExtension(filepath.Ext(name))
	if c.ETag {
		etag := fmt.Sprintf(`"%x-%x"`, fi.ModTime().Unix(), size)
		if doFresh(ctx, []byte(etag)) {
			return
		}
		util.SetHeader(ctx, cs.HeaderETag, etag)
	}

	// if file size less than 1KB, not use gzip
	// or the content is not text
	if size < minCompressSize || !textRegexp.MatchString(contentType) {
		compression = false
	}

	err = ctx.ServeContent(f, fi.Name(), fi.ModTime(), compression)
	return
}

func serveAsset(ctx iris.Context, asset *asset.Asset, filename string, c *StaticServeConfig) (err error) {
	if !asset.Exists(filename) {
		resErr(ctx, &util.HTTPError{
			StatusCode: http.StatusNotFound,
			Category:   util.ErrCategoryLogic,
			Code:       util.ErrCodeValidate,
			Message:    filename + " is not found in asset",
		})
		return
	}
	compression := c.Compression
	buf := asset.Get(filename)
	contentType := mime.TypeByExtension(filepath.Ext(filename))
	ctx.ContentType(contentType)
	size := len(buf)
	// if file size less than 1KB, not use gzip
	// or the content is not text
	if size < minCompressSize || !textRegexp.MatchString(contentType) {
		compression = false
	}
	if compression {
		data, _ := util.Gzip(buf, 0)
		// if gzip success, set content encoding
		if len(data) != 0 {
			buf = data
			util.SetHeader(ctx, cs.HeaderContentEncoding, cs.Gzip)
		}
	}
	if c.ETag {
		etag := util.GenETag(buf)
		if doFresh(ctx, []byte(etag)) {
			return
		}
		util.SetHeader(ctx, cs.HeaderETag, etag)
	}
	_, err = ctx.Write(buf)
	return
}

// StaticServe static server handler
func StaticServe(conf StaticServeConfig) iris.Handler {
	staticPath := conf.Path
	asset := conf.Asset
	maxAge := conf.MaxAge
	sMaxAge := conf.SMaxAge
	header := conf.Header
	return func(ctx iris.Context) {
		// 取url中*那部分的参数
		filename := ""

		v, ok := ctx.Params().GetEntryAt(0)
		if ok {
			filename = v.ValueRaw.(string)
		} else {
			filename = ctx.Request().URL.Query().Get("file")
		}
		if filename == "" {
			resErr(ctx, &util.HTTPError{
				StatusCode: http.StatusBadRequest,
				Category:   util.ErrCategoryLogic,
				Code:       util.ErrCodeValidate,
				Message:    "file name can not be nil",
			})
			return
		}
		for k, v := range header {
			util.SetHeader(ctx, k, v)
		}

		if maxAge != "" && sMaxAge != "" {
			util.SetCacheWithSMaxAge(ctx, maxAge, sMaxAge)
		} else if maxAge != "" {
			util.SetCache(ctx, maxAge)
		}

		var err error
		if asset != nil {
			err = serveAsset(ctx, asset, filename, &conf)
		} else {
			file := path.Join(staticPath, filename)
			err = serveFile(ctx, file, &conf)
		}
		if err != nil {
			util.GetLogger().Error("serve static file fail",
				zap.String("uri", ctx.Request().RequestURI),
				zap.Error(err),
			)
		}
	}
}
