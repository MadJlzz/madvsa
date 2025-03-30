package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Storer interface {
	Store(r io.Reader) error
}

type StorageDestination struct {
	Prefix      string
	Destination string
}

func NewStorageDestination(outputPath string, executionId string) (*StorageDestination, error) {
	splitPath := strings.Split(outputPath, "://")
	if len(splitPath) != 2 {
		return nil, fmt.Errorf("path %s should be formated like protocol://path/to/file", outputPath)
	}

	finalDestination := path.Join(splitPath[1], executionId)
	return &StorageDestination{Prefix: splitPath[0], Destination: finalDestination}, nil
}

func Store(r io.Reader, sd *StorageDestination) error {
	var s Storer
	switch sd.Prefix {
	case "file":
		s = &File{Destination: sd.Destination}
	default:
		return fmt.Errorf("prefix %s for destination %s is not supported", sd.Prefix, sd.Destination)
	}
	return s.Store(r)
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
