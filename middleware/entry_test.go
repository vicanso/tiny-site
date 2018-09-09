package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vicanso/tiny-site/util"
)

func TestNewEntry(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "http://aslant.site/", nil)
	w := httptest.NewRecorder()
	fn := NewEntry()
	ctx := util.NewContext(w, r)
	fn(ctx)
	logger := util.GetLogger()
	if logger == nil {
		t.Fatalf("entry middle should create a user logger")
	}
}
