package main

import (
	"io"
	"os"
	"path/filepath"
)

type Storer interface {
	Store(r io.Reader) error
}

type File struct {
	Destination string
}

func (f *File) Store(r io.Reader) error {
	out := filepath.Clean(f.Destination)
	fd, err := os.Create(out)
	if err != nil {
		return err
	}
	if _, err = io.Copy(fd, r); err != nil {
		return err
	}
	return nil
}
