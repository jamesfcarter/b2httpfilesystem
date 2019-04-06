package b2httpfilesystem

import (
	"net/http"
	"strings"

	blazer "github.com/kurin/blazer/b2"
)

// Filesystem implements the http.FileSystem interface for a B2 bucket
type Filesystem blazer.Bucket

func New() *Filesystem {
	return nil
}

func (f *Filesystem) Open(path string) (http.File, error) {
	b2 := (*blazer.Bucket)(f)
	path = strings.TrimPrefix(path, "/")
	return &File{bucket: b2, object: b2.Object(path)}, nil
}
