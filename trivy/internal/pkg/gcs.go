package pkg

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
		log.Fatalf("failed to gcs client: %s\n", err)
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
		return fmt.Errorf("io.Copy: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}

	fmt.Println("Successfully wrote test object")
	return nil
}
