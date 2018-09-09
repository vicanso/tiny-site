package middleware

import (
	"strconv"

	"github.com/kataras/iris"
	"github.com/vicanso/tiny-site/cs"
	"github.com/vicanso/tiny-site/util"
	"go.uber.org/zap"
)

// NewRespond 新建响应处理
func NewRespond() iris.Handler {
	logger := util.GetLogger()
	return func(ctx iris.Context) {
		ctx.Next()
		body := util.GetBody(ctx)
		if body == nil {
			return
		}
		var err error
		contentType := ctx.GetContentType()
		switch body.(type) {
		case string:
			_, err = ctx.WriteString(body.(string))
		case []byte:
			if contentType == "" {
				ctx.ContentType(cs.ContentBinaryHeaderValue)
			}
			buf := body.([]byte)
			util.SetHeader(ctx, cs.HeaderContentLength, strconv.Itoa(len(buf)))
			_, err = ctx.Write(buf)
		default:
			_, err = ctx.JSON(body)
		}
		if err != nil {
			logger.Error("response fail",
				zap.String("uri", ctx.Request().RequestURI),
				zap.Error(err),
			)
		}
	}
}
