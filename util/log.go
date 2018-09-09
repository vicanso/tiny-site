package util

import (
	"github.com/kataras/iris"
	"go.uber.org/zap"
)

var (
	defaultLogger *zap.Logger
)

func init() {
	l, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defaultLogger = l
}

const (
	// LogCategory log category field
	LogCategory = "category"
	// LogTrack log track field
	LogTrack = "track"

	// LogAccess access log category
	LogAccess = "access"
	// LogTracker tracker log category
	LogTracker = "tracker"
	// LogUser user log category
	LogUser = "user"
)

// GetLogger get logger
func GetLogger() *zap.Logger {
	return defaultLogger
}

// CreateAccessLogger 创建access logger
func CreateAccessLogger() *zap.Logger {
	return defaultLogger.With(zap.String(LogCategory, LogAccess))
}

// CreateTrackerLogger 创建tracker logger
func CreateTrackerLogger() *zap.Logger {
	return defaultLogger.With(zap.String(LogCategory, LogTracker))
}

// CreateUserLogger 创建user logger
func CreateUserLogger(ctx iris.Context) *zap.Logger {
	return defaultLogger.With(
		zap.String(LogCategory, LogUser),
		zap.String(LogTrack, GetTrackID(ctx)))
}

// SetContextLogger 设置logger
func SetContextLogger(ctx iris.Context, logger *zap.Logger) {
	ctx.Values().Set(Logger, logger)
}

// GetContextLogger 获取logger
func GetContextLogger(ctx iris.Context) *zap.Logger {
	logger := ctx.Values().Get(Logger)
	if logger == nil {
		return nil
	}
	return logger.(*zap.Logger)
}
