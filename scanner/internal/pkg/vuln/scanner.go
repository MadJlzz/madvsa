package vuln

import (
	"bytes"
	"context"
)

type Scanner interface {
	Scan(ctx context.Context, image string) (*bytes.Buffer, error)
}
