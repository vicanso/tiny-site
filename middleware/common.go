package middleware

import (
	"time"

	"github.com/vicanso/tiny-site/cs"
	"github.com/vicanso/tiny-site/model"

	"github.com/kataras/iris"
	"github.com/vicanso/tiny-site/util"
)

// IsLogined check login status，if not, will return error
func IsLogined(ctx iris.Context) {
	if util.GetAccount(ctx) == "" {
		resErr(ctx, util.ErrNeedLogined)
		return
	}
	ctx.Next()
}

// IsAnonymous check login status, if yes, will return error
func IsAnonymous(ctx iris.Context) {
	if util.GetAccount(ctx) != "" {
		resErr(ctx, util.ErrLoginedAlready)
		return
	}
	ctx.Next()
}

// WaitFor at least wait for duration
func WaitFor(d time.Duration) iris.Handler {
	ns := d.Nanoseconds()
	return func(ctx iris.Context) {
		start := time.Now()
		ctx.Next()
		use := time.Now().UnixNano() - start.UnixNano()
		if use < ns {
			time.Sleep(time.Duration(ns-use) * time.Nanosecond)
		}
	}
}

// IsSu check the user roles include su
func IsSu(ctx iris.Context) {
	account := util.GetAccount(ctx)
	if account == "" {
		resErr(ctx, util.ErrNeedLogined)
		return
	}
	sess := util.GetSession(ctx)
	roles := sess.GetStringSlice(cs.SessionRolesField)
	if !util.ContainsString(roles, model.UserRoleSu) {
		resErr(ctx, util.ErrUserForbidden)
		return
	}
	ctx.Next()
}

// IsNilQuery check the query is nil
func IsNilQuery(ctx iris.Context) {
	if ctx.Request().URL.RawQuery != "" {
		resErr(ctx, util.ErrQueryShouldBeNil)
		return
	}
	ctx.Next()
	return
}
