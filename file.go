package b2httpfilesystem

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"os"

	blazer "github.com/kurin/blazer/b2"
)

type File struct {
	bucket  *blazer.Bucket
	object  *blazer.Object
	content *bytes.Reader
	closed  bool
}

func (f *File) download() error {
	if f.closed {
		return os.ErrClosed
	}
	if f.content != nil {
		return nil
	}
	log.Printf("B2 downloading %s\n", f.object.Name())
	reader := f.object.NewReader(context.Background())
	defer reader.Close()
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	f.content = bytes.NewReader(data)
	return nil
}

func (f *File) Read(p []byte) (int, error) {
	err := f.download()
	if err != nil {
		return 0, err
	}
	return f.content.Read(p)
}

func (f *File) Close() error {
	if f.closed {
		return os.ErrClosed
	}
	f.content = nil
	f.closed = true
	return nil
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	err := f.download()
	if err != nil {
		return 0, err
	}
	return f.content.Seek(offset, whence)
}

func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	var result []os.FileInfo
	iter := f.bucket.List(context.Background(),
		blazer.ListDelimiter("/"),
		blazer.ListPrefix(f.object.Name()))
	for iter.Next() {
		attrs, err := iter.Object().Attrs(context.Background())
		if err != nil {
			return nil, err
		}
		result = append(result, &info{attrs: attrs})
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (f *File) Stat() (os.FileInfo, error) {
	attrs, err := f.object.Attrs(context.Background())
	if err != nil {
		return nil, err
	}
	return &info{attrs: attrs}, nil
}
