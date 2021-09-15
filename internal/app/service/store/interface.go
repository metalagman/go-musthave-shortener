package store

import (
	"errors"
	"fmt"
)

var (
	ErrBadInput   = errors.New("bad input")
	ErrEmptyInput = fmt.Errorf("%w: empty url", ErrBadInput)
	ErrNotFound   = errors.New("not found")
)

// Store of the url data
type Store interface {
	// WriteURL to storage, returns short URL
	WriteURL(url string, uid string) (string, error)
	// ReadURL from storage
	ReadURL(id string) (string, error)
	// ReadUserURLs from db
	ReadUserURLs(uid string) []StoredURL
}

type StoredURL struct {
	ID          string
	ShortURL    string
	OriginalURL string
}
