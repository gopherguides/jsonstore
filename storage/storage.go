package storage

import "io"

type Storer interface {
	Reader(path string) (io.ReadCloser, error)
	Writer(path string) (io.WriteCloser, error)
}
