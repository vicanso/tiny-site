package controller

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/vicanso/session"
	"github.com/vicanso/tiny-site/util"
)

func newfileUploadRequest(uri string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req := httptest.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func TestFileCtrl(t *testing.T) {
	ctrl := fileCtrl{}
	var uploadInfo *uploadInfoResponse
	t.Run("upload", func(t *testing.T) {
		r, err := newfileUploadRequest("http://127.0.0.1/", "file", "../assets/ai.jpeg")
		if err != nil {
			t.Fatalf("create upload request fail, %v", err)
		}
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, r)
		ctrl.upload(ctx)
		if ctx.GetStatusCode() != http.StatusCreated {
			t.Fatalf("upload file fail")
		}
		uploadInfo = util.GetBody(ctx).(*uploadInfoResponse)
		if uploadInfo.ID == "" || uploadInfo.FileType != "jpeg" {
			t.Fatalf("upload file info is wrong")
		}
	})

	t.Run("save", func(t *testing.T) {
		buf := []byte(`{
			"file": "` + uploadInfo.ID + `",
			"category": "test",
			"fileType": "` + uploadInfo.FileType + `",
			"maxAge": "1h"
		}`)
		r := httptest.NewRequest(http.MethodPost, "http://127.0.0.1/", nil)
		w := httptest.NewRecorder()

		ctx := util.NewContext(w, r)
		sess := session.Mock(session.M{
			"fetched": true,
			"data": session.M{
				"account": "vicanso",
			},
		})
		util.SetSession(ctx, sess)
		util.SetRequestBody(ctx, buf)
		ctrl.save(ctx)
		if ctx.GetStatusCode() != http.StatusCreated {
			t.Fatalf("save file fail")
		}
	})

	t.Run("save file is expired", func(t *testing.T) {
		buf := []byte(`{
			"id": "01CQ0YRSERJB95SNBNF2VBNGN5",
			"category": "test",
			"fileType": "jpeg",
			"maxAge": "1h"
		}`)
		r := httptest.NewRequest(http.MethodPost, "http://127.0.0.1/", nil)
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, r)
		util.SetRequestBody(ctx, buf)
		ctrl.save(ctx)
		if ctx.GetStatusCode() != http.StatusBadRequest {
			t.Fatalf("file expired should return error")
		}
	})

	t.Run("save file with wrong id", func(t *testing.T) {
		buf := []byte(`{
			"id": "01CQ0YRSERJB9",
			"category": "test",
			"fileType": "jpeg",
			"maxAge": "1h"
		}`)
		r := httptest.NewRequest(http.MethodPost, "http://127.0.0.1/", nil)
		w := httptest.NewRecorder()
		ctx := util.NewContext(w, r)
		util.SetRequestBody(ctx, buf)
		ctrl.save(ctx)
		if ctx.GetStatusCode() != http.StatusBadRequest {
			t.Fatalf("wrong params should return error")
		}
	})
}
