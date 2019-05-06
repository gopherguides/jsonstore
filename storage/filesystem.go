package storage

import (
	"io"
	"os"
)

type FileSystem struct{}

func NewFileSystem() Storer {
	return FileSystem{}
}

func (fs FileSystem) Reader(path string) (io.ReadCloser, error) {
	f, err := os.Open(path)
	if err != nil {
		// If we don't have a file, that's ok
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return f, nil
}

func (fs FileSystem) Writer(path string) (io.WriteCloser, error) {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return f, nil
}
