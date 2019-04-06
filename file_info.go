package b2httpfilesystem

import (
	"os"
	"path/filepath"
	"syscall"
	"time"

	blazer "github.com/kurin/blazer/b2"
)

type info struct {
	attrs *blazer.Attrs
}

func (i *info) Name() string {
	return filepath.Base(i.attrs.Name)
}

func (i *info) Size() int64 {
	return i.attrs.Size
}

func (i *info) Mode() os.FileMode {
	if i.IsDir() {
		return os.ModeDir
	}
	return os.FileMode(0644)
}

func (i *info) ModTime() time.Time {
	return i.attrs.LastModified
}

func (i *info) IsDir() bool {
	return i.attrs.Status == blazer.Folder
}

func (i *info) Sys() interface{} {
	return &syscall.Stat_t{}
}
