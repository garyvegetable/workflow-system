package storage

import "io"

type Storage interface {
	Upload(key string, reader io.Reader, size int64, contentType string) error
	Download(key string) (io.ReadCloser, error)
	GetURL(key string) (string, error)
	Delete(key string) error
}
