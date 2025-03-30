package pkg

import (
	"context"
	"io"
	"net/url"
	"os"
	"path/filepath"
)

func init() {
	DefaultStoreFactory().Register("file", &FileStorage{})
}

type FileStorage struct{}

func (f *FileStorage) Store(_ context.Context, r io.Reader, destination *url.URL) error {
	out := filepath.Clean(destination.Path)
	fd, err := os.Create(out)
	if err != nil {
		return err
	}
	if _, err = io.Copy(fd, r); err != nil {
		return err
	}
	return nil
}
