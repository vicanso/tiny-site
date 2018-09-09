package middleware

import (
	"net/http"
	"regexp"

	"go.uber.org/zap"

	"github.com/kataras/iris"
	"github.com/vicanso/tiny-site/util"
)

const (
	// HandleSuccess handle success
	HandleSuccess = iota + 1
	// HandleFail handle fail
	HandleFail
)

var (
	trackerLogger = util.CreateTrackerLogger()
	// defaultMaskFields 需要替换为***的字段
	defaultMaskFields = regexp.MustCompile(`password`)
	defaultOnTrack    = func(info *TrackerInfo) {
		trackerLogger.Info("",
			zap.String("tracker", info.Category),
			zap.String("trackId", info.TrackID),
			zap.String("account", info.Account),
			zap.Any("result", info.Result),
			zap.Any("query", info.Query),
			zap.Any("params", info.Params),
			zap.Any("form", info.Form),
			zap.Any("body", info.Body),
		)
	}
)

type (
	// OnTrack on track function
	OnTrack func(*TrackerInfo)
	// TrackerInfo the info of tracker
	TrackerInfo struct {
		Category string            `json:"category"`
		TrackID  string            `json:"trackId,omitempty"`
		Account  string            `json:"account,omitempty"`
		Query    map[string]string `json:"query,omitempty"`
		Params   iris.Map          `json:"params,omitempty"`
		Form     iris.Map          `json:"form,omitempty"`
		Result   int               `json:"result,omitempty"`
		Body     interface{}       `json:"body,omitempty"`
	}
	// TrackerConfig tracker config
	TrackerConfig struct {
		Query    bool
		Params   bool
		Form     bool
		Response bool
		Mask     *regexp.Regexp
		OnTrack  OnTrack
	}
)

// converQuery convert query to json string
func converQuery(ctx iris.Context, reg *regexp.Regexp) map[string]string {
	query := util.GetRequestQuery(ctx)
	m := make(map[string]string)
	for k, v := range query {
		if reg != nil && reg.MatchString(k) {
			m[k] = "***"
			continue
		}
		m[k] = v
	}
	return m
}

// convertRequestBody convert post body to string
func convertRequestBody(ctx iris.Context, reg *regexp.Regexp) iris.Map {
	m := iris.Map{}
	json.Unmarshal(util.GetRequestBody(ctx), &m)
	for k := range m {
		if reg != nil && reg.MatchString(k) {
			m[k] = "***"
		}
	}
	return m
}

// convertParams convert params to string
func convertParams(ctx iris.Context, reg *regexp.Regexp) iris.Map {
	m := iris.Map{}
	ctx.Params().Visit(func(k, v string) {
		if reg != nil && reg.MatchString(k) {
			m[k] = "***"
		} else {
			m[k] = v
		}
	})
	return m
}

// NewTracker new a tracker
func NewTracker(category string, conf TrackerConfig) iris.Handler {
	mask := conf.Mask
	return func(ctx iris.Context) {
		info := &TrackerInfo{
			Category: category,
			TrackID:  getTrackID(ctx),
		}
		if conf.Query {
			info.Query = converQuery(ctx, mask)
		}
		if conf.Form {
			info.Form = convertRequestBody(ctx, mask)
		}
		if conf.Params {
			info.Params = convertParams(ctx, mask)
		}
		ctx.Next()
		// get account should after fetch session
		info.Account = getAccount(ctx)
		status := ctx.GetStatusCode()
		if conf.Response {
			info.Body = util.GetBody(ctx)
		}
		if status < http.StatusOK || status >= http.StatusBadRequest {
			info.Result = HandleFail
		} else {
			info.Result = HandleSuccess
		}
		if conf.OnTrack != nil {
			conf.OnTrack(info)
		}
	}
}

// NewDefaultTracker new a tracker with default config
func NewDefaultTracker(category string, onTrack OnTrack) iris.Handler {
	// defaultTrackerConfig default config
	if onTrack == nil {
		onTrack = defaultOnTrack
	}
	config := TrackerConfig{
		Query:    true,
		Params:   true,
		Form:     true,
		Response: true,
		Mask:     defaultMaskFields,
		OnTrack:  onTrack,
	}

	return NewTracker(category, config)
}
