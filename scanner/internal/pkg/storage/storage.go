package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
)

const (
	FileScheme   = "file"
	GoogleScheme = "gcs"
)

type Storer interface {
	Store(ctx context.Context, r io.Reader, destination *url.URL) error
}

type StorerFactory struct {
	s Storer
}

func NewStorerFactory(ctx context.Context, url *url.URL) (*StorerFactory, error) {
	var sf StorerFactory
	switch url.Scheme {
	case FileScheme:
		sf.s = &FileStorage{}
	case GoogleScheme:
		sf.s = NewGoogleBlobStorage(ctx)
	default:
		return nil, fmt.Errorf("unsupported storage scheme: %s", url.Scheme)
	}
	return &sf, nil
}

func (sf *StorerFactory) Store(ctx context.Context, r io.Reader, destination *url.URL) error {
	return sf.s.Store(ctx, r, destination)
}
