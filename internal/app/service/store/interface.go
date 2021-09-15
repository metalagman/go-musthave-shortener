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

// Store of the db served by store
type Store interface {
	// WriteURL to storage, returns short URL
	WriteURL(url string) (string, error)
	// ReadURL from storage
	ReadURL(id string) (string, error)
}
