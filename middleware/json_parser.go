package middleware

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/vicanso/tiny-site/util"

	"github.com/kataras/iris"
)

const (
	// 默认为50kb
	defaultRequestJSONLimit = 50 * 1024
)

type (
	// JSONParserConfig config配置
	JSONParserConfig struct {
		// 数据大小限制
		Limit int
	}
)

// NewJSONParser 创建新的json parser handler
func NewJSONParser(conf JSONParserConfig) iris.Handler {
	limit := defaultRequestJSONLimit
	if conf.Limit != 0 {
		limit = conf.Limit
	}
	return func(ctx iris.Context) {
		method := ctx.Method()
		if method != http.MethodPost && method != http.MethodPatch && method != http.MethodPut {
			ctx.Next()
			return
		}
		req := ctx.Request()
		contentType := req.Header.Get("Content-Type")

		if !strings.HasPrefix(contentType, "application/json") {
			ctx.Next()
			return
		}
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			resErr(ctx, err)
			return
		}
		if limit != 0 && len(body) > limit {
			resErr(ctx, util.ErrRequestJSONTooLarge)
			return
		}
		util.SetRequestBody(ctx, body)
		ctx.Next()
	}
}
