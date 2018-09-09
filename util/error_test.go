package util

import (
	"errors"
	"testing"
)

func TestNewHTTPError(t *testing.T) {
	statusCode := 500
	message := "出错了"
	code := "custom-error-cdoe"
	he := NewHTTPError(statusCode, message, code)
	if he.StatusCode != statusCode {
		t.Fatalf("new error with status code fail")
	}
	if he.Message != message {
		t.Fatalf("new error with message fail")
	}
	if he.Code != code {
		t.Fatalf("new error with code fail")
	}
}

func TestNewJSONParseError(t *testing.T) {
	err := errors.New("json parse error")
	he := NewJSONParseError(err)
	if he.Category != ErrCategoryJSON {
		t.Fatalf("json parse error category is wrong")
	}
	if he.Message != err.Error() {
		t.Fatalf("json parse error message is wrong")
	}
}

func TestNewValidateError(t *testing.T) {
	err := errors.New("validate error")
	he := NewValidateError(err)
	if he.Category != ErrCategoryValidate {
		t.Fatalf("validate error category is wrong")
	}
	if he.Message != err.Error() {
		t.Fatalf("validate error message is wrong")
	}
}
