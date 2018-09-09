package util

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/h2non/gock"
)

func TestCreateRequestError(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/", nil)
	w := &http.Response{
		StatusCode: http.StatusBadRequest,
	}
	message := "出错了"
	err := createRequestError(message, r, w)
	if err.Category != ErrCategoryRequset ||
		err.StatusCode != http.StatusBadRequest ||
		err.Message != message {
		t.Fatalf("create error fail")
	}
}

func TestHTTPGet(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution
	m := map[string]string{"foo": "bar"}
	gock.New("http://aslant.site").
		Get("/bar").
		Reply(200).
		JSON(m)

	data, err := HTTPGet("http://aslant.site/bar", map[string]string{
		"type": "1",
	})
	if err != nil {
		t.Fatalf("http get fail")
	}
	buf, _ := json.Marshal(m)
	if !bytes.Equal(buf, bytes.TrimSpace(data)) {
		t.Fatalf("http get data fail")
	}
}

func TestHTTPPost(t *testing.T) {
	defer gock.Off()
	m := map[string]string{
		"message": "出错了",
	}
	gock.New("http://aslant.site").
		Post("/login").
		Reply(500).
		JSON(m)
	data := map[string]interface{}{
		"account": "vicanso",
	}
	params := map[string]string{
		"type": "1",
	}
	_, err := HTTPPost("http://aslant.site/login?type=1", data, params)
	if err == nil {
		t.Fatalf("http post should return error")
	}
	he, ok := err.(*HTTPError)
	if !ok {
		t.Fatalf("http post error should be HTTPError")
	}
	if he.Category != ErrCategoryRequset {
		t.Fatalf("http requset error category is wrong")
	}
}

func TestHookConverError(t *testing.T) {
	h := &hookConverError{}
	t.Run("after request with error", func(t *testing.T) {

		r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/", nil)
		err := errors.New("request fail")
		_, newErr := h.AfterRequest(r, nil, err)
		he := newErr.(*HTTPError)
		if he.StatusCode != http.StatusInternalServerError ||
			he.Message != err.Error() {
			t.Fatalf("create new error fail")
		}
	})

	t.Run("after request, read response body fail", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/", nil)
		err := errors.New("read error")
		w := &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       NewErrorReadCloser(err),
		}
		_, newErr := h.AfterRequest(r, w, nil)
		if newErr.Error() != err.Error() {
			t.Fatalf("read error fail")
		}
	})

	t.Run("after request, read response but not json", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/", nil)
		w := &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       ioutil.NopCloser(strings.NewReader("abc")),
		}
		_, newErr := h.AfterRequest(r, w, nil)
		if newErr == nil {
			t.Fatalf("not json response should return error")
		}
	})
}
