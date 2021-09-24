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

// Ping allows you to perform store health check
type Ping interface {
	// Ping underlying storage and return error if it is not available
	Ping() error
}

// Store of the url data
type Store interface {
	Reader
	UserWriter
	UserReader
}

type Reader interface {
	// ReadURL from storage using provided id
	ReadURL(id string) (string, error)
}

type Writer interface {
	// WriteURL from storage using provided id
	WriteURL(id string) (string, error)
}

type UserReader interface {
	// ReadUserURLs from db
	ReadUserURLs(uid string) []StoredURL
}

type UserWriter interface {
	// WriteUserURL to storage, returns short URL
	WriteUserURL(url string, uid string) (string, error)
}

type StoredURL struct {
	ID          string
	ShortURL    string
	OriginalURL string
}
