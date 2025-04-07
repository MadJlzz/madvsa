package storage

import (
	"context"
	"fmt"
	"io"
	"log/slog"
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
	logger *slog.Logger
	s      Storer
}

func NewStorerFactory(ctx context.Context, logger *slog.Logger, url *url.URL) (*StorerFactory, error) {
	var sf StorerFactory
	sf.logger = logger
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
	sf.logger.Info("storing image", "storer", typeof(sf.s), "url", destination)
	return sf.s.Store(ctx, r, destination)
}

// Might consider putting that somewhere else if ever reused.
func typeof(v any) string {
	return fmt.Sprintf("%T", v)
}
