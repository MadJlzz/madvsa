package pkg

import (
	"context"
	"fmt"
	"io"
	"net/url"
)

type Storer interface {
	Store(ctx context.Context, r io.Reader, destination *url.URL) error
}

type StoreFactory struct {
	sm map[string]Storer
}

var defaultStoreFactory = &StoreFactory{sm: map[string]Storer{}}

func DefaultStoreFactory() *StoreFactory {
	return defaultStoreFactory
}

func (sf *StoreFactory) Register(scheme string, storer Storer) {
	sf.sm[scheme] = storer
}

func (sf *StoreFactory) GetStorer(url *url.URL) (Storer, error) {
	s, ok := sf.sm[url.Scheme]
	if !ok {
		return nil, fmt.Errorf("unsupported scheme %s", url.Scheme)
	}
	return s, nil
}
