package util

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLogger(t *testing.T) {
	if GetLogger() == nil {
		t.Fatalf("get logger fail")
	}

	if CreateAccessLogger() == nil {
		t.Fatalf("create sugger logger fail")
	}

	if CreateTrackerLogger() == nil {
		t.Fatalf("create tracker logger fail")
	}

	r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/", nil)
	ctx := NewContext(nil, r)
	if CreateUserLogger(ctx) == nil {
		t.Fatalf("create user logger fail")
	}

	logger := GetLogger()
	SetContextLogger(ctx, logger)
	if GetContextLogger(ctx) == nil {
		t.Fatalf("get context logger fail")
	}

}
