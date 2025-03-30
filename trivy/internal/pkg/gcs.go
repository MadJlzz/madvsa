package pkg

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"
	"sync"
)

func init() {
	DefaultStoreFactory().Register("gs", &GCS{})
}

type GCS struct {
	init sync.Once
	cli  *storage.Client
}

func (s *GCS) Store(ctx context.Context, r io.Reader, destination *url.URL) error {
	s.init.Do(func() {
		cli, err := storage.NewClient(ctx)
		if err != nil {
			log.Fatalf("failed to create GCS client: %s", err)
		}
		s.cli = cli
	})

	bh := s.cli.Bucket(destination.Host)
	oh := bh.Object(strings.TrimPrefix(destination.Path, "/"))

	wc := oh.NewWriter(ctx)
	if _, err := io.Copy(wc, r); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}

	fmt.Println("Successfully wrote test object")
	return nil
}
