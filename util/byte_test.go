package util

import (
	"bytes"
	"testing"
)

func TestGzip(t *testing.T) {
	buf := []byte("ABCD")
	data, err := Gzip(buf, 0)
	if err != nil {
		t.Fatalf("gzip fail, %v", err)
	}
	gzipData := []byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 114, 116, 114, 118, 1, 4, 0, 0, 255, 255, 165, 32, 23, 219, 4, 0, 0, 0}
	if !bytes.Equal(data, gzipData) {
		t.Fatalf("gzip fail")
	}
}

func TestGenETag(t *testing.T) {
	buf := []byte("")
	if GenETag(buf) != `"0-2jmj7l5rSw0yVb_vlWAYkK_YBwk="` {
		t.Fatalf("gen empty byte's etag fail")
	}
	buf = []byte("ABCD")
	if GenETag(buf) != `"4--y-FyIVn88jOm3mcfFRkLQx7QfY="` {
		t.Fatalf("gen etag fail")
	}
}
