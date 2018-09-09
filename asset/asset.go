package asset

import (
	"net/http"

	"github.com/gobuffalo/packr"
)

var box packr.Box

type (
	// Asset static asset
	Asset struct {
	}
)

func init() {
	box = packr.NewBox("../assets")
}

// New create an asset instance
func New() *Asset {
	return &Asset{}
}

// Open open the file
func (a *Asset) Open(filename string) (http.File, error) {
	return box.Open(filename)
}

// Get the the data of file
func (a *Asset) Get(filename string) []byte {
	return box.Bytes(filename)
}

// Exists check the file exists
func (a *Asset) Exists(filename string) bool {
	return box.Has(filename)
}
