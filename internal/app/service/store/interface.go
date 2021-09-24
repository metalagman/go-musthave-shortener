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
	Writer
	UserDataReader
}

type Reader interface {
	// ReadURL from storage using provided id
	ReadURL(id string) (string, error)
}

type UserDataReader interface {
	// ReadAllURLs from db
	ReadAllURLs(uid string) []Record
}

type Writer interface {
	// WriteURL to storage, returns short Record
	WriteURL(url string, uid string) (string, error)
}

type BatchWriter interface {
	BatchWrite(uid string, in []Record) ([]Record, error)
}

type RecordID string

type Record struct {
	ID            string
	ShortURL      string
	OriginalURL   string
	CorrelationID string
}
