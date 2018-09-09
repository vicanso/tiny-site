package asset

import (
	"testing"
)

func TestAsset(t *testing.T) {
	as := New()
	filename := "index.html"
	t.Run("open", func(t *testing.T) {
		f, err := as.Open(filename)
		if err != nil {
			t.Fatalf("open fail, %v", err)
		}
		fi, err := f.Stat()
		if err != nil {
			t.Fatalf("stat fail, %v", err)
		}
		if fi.Name() != filename {
			t.Fatalf("get stat info fail")
		}
	})

	t.Run("get", func(t *testing.T) {
		buf := as.Get(filename)
		if len(buf) == 0 {
			t.Fatalf("get file data fail")
		}
	})

	t.Run("exists", func(t *testing.T) {
		if !as.Exists(filename) {
			t.Fatalf("check exists fail")
		}
	})
}
