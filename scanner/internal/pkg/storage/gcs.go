package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"
)

type GoogleBlobStorage struct {
	cli *storage.Client
}

func NewGoogleBlobStorage(ctx context.Context) *GoogleBlobStorage {
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("gcs: failed to create storage client: %s\n", err)
	}
	return &GoogleBlobStorage{
		cli: storageClient,
	}
}

func (gcs *GoogleBlobStorage) Store(ctx context.Context, r io.Reader, destination *url.URL) error {
	bh := gcs.cli.Bucket(destination.Host)
	oh := bh.Object(strings.TrimPrefix(destination.Path, "/"))

	wc := oh.NewWriter(ctx)
	if _, err := io.Copy(wc, r); err != nil {
		return fmt.Errorf("gcs: io.Copy: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("gcs: Writer.Close: %w", err)
	}

	fmt.Println("Successfully wrote test object")
	return nil
}
