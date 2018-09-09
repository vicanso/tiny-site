package util

import (
	"io"
	"net/http"

	"github.com/kataras/iris"
)

const (
	// ErrCodeJSON josn parse error	code
	ErrCodeJSON = "101"
	// ErrCodeValidate validate error code
	ErrCodeValidate = "102"
	// ErrCodeSessionFetch session fetch error code
	ErrCodeSessionFetch = "103"
	// ErrCodeSessionCommit session commit error code
	ErrCodeSessionCommit = "104"
	// ErrCodeUser user error code
	ErrCodeUser = "110"
	// ErrCodeFile file error code
	ErrCodeFile = "111"
)

const (
	// ErrCategoryJSON json出错的类别
	ErrCategoryJSON = "json"
	// ErrCategoryValidate validate出错的类别
	ErrCategoryValidate = "validate"
	// ErrCategoryRequset request出错类别
	ErrCategoryRequset = "request"
	// ErrCategoryLogic 逻辑相关校验出错
	ErrCategoryLogic = "logic"
	// ErrCategorySession session出错
	ErrCategorySession = "session"
)

var (
	// ErrRequestJSONTooLarge request too large
	ErrRequestJSONTooLarge = &HTTPError{
		StatusCode: http.StatusRequestEntityTooLarge,
		Message:    "request post json too large",
	}
	// ErrTooManyRequest too many requset
	ErrTooManyRequest = &HTTPError{
		StatusCode: http.StatusTooManyRequests,
		Message:    "too many request",
	}
	// ErrServiceUnavailable service unavailable
	ErrServiceUnavailable = &HTTPError{
		StatusCode: http.StatusServiceUnavailable,
		Message:    "service unavailable",
	}
	// ErrQueryShouldBeNil query should be nil
	ErrQueryShouldBeNil = &HTTPError{
		StatusCode: http.StatusBadRequest,
		Category:   ErrCategoryLogic,
		Message:    "query should be nil",
	}
	// ErrNoTrackKey no track key
	ErrNoTrackKey = &HTTPError{
		StatusCode: http.StatusBadRequest,
		Category:   ErrCategoryLogic,
		Message:    "track key is not found",
	}
	// ErrLoginedAlready logined already
	ErrLoginedAlready = &HTTPError{
		StatusCode: http.StatusBadRequest,
		Category:   ErrCategoryLogic,
		Code:       ErrCodeUser,
		Message:    "user is logined, please logout first",
	}
	// ErrNeedLogined login first
	ErrNeedLogined = &HTTPError{
		StatusCode: http.StatusUnauthorized,
		Category:   ErrCategoryLogic,
		Code:       ErrCodeUser,
		Message:    "please login first",
	}
	// ErrAccountOrPasswordWrong account or password is wrong
	ErrAccountOrPasswordWrong = &HTTPError{
		StatusCode: http.StatusBadRequest,
		Category:   ErrCategoryLogic,
		Code:       ErrCodeUser,
		Message:    "account or password is wrong",
	}
	// ErrLoginTokenNil login token can not be nil
	ErrLoginTokenNil = &HTTPError{
		StatusCode: http.StatusBadRequest,
		Category:   ErrCategoryLogic,
		Code:       ErrCodeUser,
		Message:    "login token can not be nil",
	}
	// ErrUserForbidden forbidden the function
	ErrUserForbidden = &HTTPError{
		StatusCode: http.StatusForbidden,
		Category:   ErrCategoryLogic,
		Code:       ErrCodeUser,
		Message:    "this function is forbidden",
	}
	// ErrFileNameIsNil file name can not be nil
	ErrFileNameIsNil = &HTTPError{
		StatusCode: http.StatusBadRequest,
		Category:   ErrCategoryLogic,
		Code:       ErrCodeValidate,
		Message:    "file name can not be nil",
	}
)

type (
	// HTTPError http error
	HTTPError struct {
		StatusCode int
		Category   string
		Message    string
		Code       string // custom code
		Internal   error  // Stores the error returned by an external dependency
		Extra      iris.Map
	}
	errReadCloser struct {
		customErr error
	}
)

// Error makes it compatible with `error` interface.
func (he *HTTPError) Error() string {
	return he.Message
}

// Read read function
func (er *errReadCloser) Read(p []byte) (n int, err error) {
	return 0, er.customErr
}

// Close close function
func (er *errReadCloser) Close() error {
	return nil
}

// NewHTTPError creates a new HTTPError instance.
func NewHTTPError(statusCode int, args ...string) *HTTPError {
	he := &HTTPError{
		StatusCode: statusCode,
	}
	if len(args) > 0 {
		he.Message = args[0]
		if len(args) > 1 {
			he.Code = args[1]
		}
	}
	return he
}

// NewJSONParseError 创建新的json parse error
func NewJSONParseError(err error) *HTTPError {
	he := NewHTTPError(http.StatusBadRequest, err.Error(), ErrCodeJSON)
	he.Category = ErrCategoryJSON
	return he
}

// NewValidateError 创新新的validate error
func NewValidateError(err error) *HTTPError {
	he := NewHTTPError(http.StatusBadRequest, err.Error(), ErrCodeValidate)
	he.Category = ErrCategoryValidate
	return he
}

// NewErrorReadCloser 创建一个出错的reader
func NewErrorReadCloser(err error) io.ReadCloser {
	r := &errReadCloser{
		customErr: err,
	}
	return r
}
