package middleware

import (
	"net/http"
	"time"

	"github.com/vicanso/tiny-site/util"

	"github.com/go-redis/redis"
	"github.com/kataras/iris"
	"github.com/vicanso/session"
)

const defaultMemoryStoreSize = 1024

type (
	// SessionConfig session's config
	SessionConfig struct {
		Cookie       string
		CookieMaxAge time.Duration
		CookiePath   string
		Expires      time.Duration
		Keys         []string
	}
)

// NewSession 创建新的session中间件
func NewSession(client *redis.Client, conf SessionConfig) iris.Handler {
	var store session.Store
	if client != nil {
		store = session.NewRedisStore(client, nil)
	} else {
		store, _ = session.NewMemoryStore(defaultMemoryStoreSize)
	}
	opts := &session.Options{
		Store:        store,
		Key:          conf.Cookie,
		MaxAge:       int(conf.Expires.Seconds()),
		CookieKeys:   conf.Keys,
		CookieMaxAge: int(conf.CookieMaxAge.Seconds()),
		CookiePath:   conf.CookiePath,
	}
	return func(ctx iris.Context) {
		res := ctx.ResponseWriter()
		req := ctx.Request()
		sess := session.New(req, res, opts)
		_, err := sess.Fetch()
		if err != nil {
			resErr(ctx, &util.HTTPError{
				StatusCode: http.StatusInternalServerError,
				Category:   util.ErrCategorySession,
				Message:    err.Error(),
				Code:       util.ErrCodeSessionFetch,
			})
			return
		}
		util.SetSession(ctx, sess)
		ctx.Next()
		err = sess.Commit()
		if err != nil {
			resErr(ctx, &util.HTTPError{
				StatusCode: http.StatusInternalServerError,
				Category:   util.ErrCategorySession,
				Message:    err.Error(),
				Code:       util.ErrCodeSessionCommit,
			})
			return
		}
	}
}
