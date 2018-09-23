package controller

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/vicanso/session"
	"github.com/vicanso/tiny-site/util"
)

const pngData = "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"

func newfileUploadRequest(uri string, paramName, path string) (*http.Request, error) {

	buf, err := base64.StdEncoding.DecodeString(pngData)
	if err != nil {
		return nil, err
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, bytes.NewReader(buf))

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
		r, err := newfileUploadRequest("http://127.0.0.1/", "file", "../assets/ai.png")
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
		if uploadInfo == nil || uploadInfo.ID == "" || uploadInfo.FileType != "png" {
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
