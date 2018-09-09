package middleware

import (
	"github.com/kataras/iris"
	"github.com/vicanso/tiny-site/util"
)

// NewEntry create a new entry
func NewEntry() iris.Handler {
	return func(ctx iris.Context) {
		util.SetNoCache(ctx)
		logger := util.CreateUserLogger(ctx)
		util.SetContextLogger(ctx, logger)
		ctx.Next()
	}
}
