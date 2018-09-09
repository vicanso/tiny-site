package middleware

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/vicanso/tiny-site/util"
)

var (
	json       = jsoniter.ConfigCompatibleWithStandardLibrary
	getTrackID = util.GetTrackID
	getAccount = util.GetAccount
	resErr     = util.ResErr
)
