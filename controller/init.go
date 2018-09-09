package controller

import (
	"github.com/vicanso/tiny-site/middleware"
	"github.com/vicanso/tiny-site/util"
)

var (
	res             = util.Res
	resNoContent    = util.ResNoContent
	resCreated      = util.ResCreated
	resJPEG         = util.ResJPEG
	resPNG          = util.ResPNG
	resWEBP         = util.ResWEBP
	resErr          = util.ResErr
	setCache        = util.SetCache
	setNoCache      = util.SetNoCache
	setNoStore      = util.SetNoStore
	validate        = util.Validate
	getRequestBody  = util.GetRequestBody
	getRequestQuery = util.GetRequestQuery
	getUTCNow       = util.GetUTCNow
	getNow          = util.GetNow

	newDefaultTracker = middleware.NewDefaultTracker
	newTracker        = middleware.NewTracker
)
