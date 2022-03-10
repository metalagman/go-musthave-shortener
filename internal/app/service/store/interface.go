//go:generate mockgen -source=./interface.go -destination=./mock/store.go -package=storemock
/*
Package store provides interfaces for storing data.
*/
package store

import (
	"errors"
	"fmt"
)

var (
	ErrBadInput   = errors.New("bad input")
	ErrEmptyInput = fmt.Errorf("empty url: %w", ErrBadInput)
	ErrNotFound   = errors.New("not found")
	ErrDeleted    = errors.New("deleted")
	ErrConflict   = &ConflictError{}
)

type ConflictError struct {
	ExistingURL string
}

func (e ConflictError) Error() string {
	return "conflict"
}

// HealthChecker allows you to perform store health check
type HealthChecker interface {
	// HealthCheck underlying storage and return error if it is not available
	HealthCheck() error
}

// Store of the url data
type Store interface {
	Reader
	Writer
	UserDataReader
}

// Reader allows you to read short urls.
type Reader interface {
	// ReadURL from storage using provided id.
	ReadURL(id string) (string, error)
}

// UserDataReader allows you to read user short urls.
type UserDataReader interface {
	// ReadUserData from db
	ReadUserData(uid string) []Record
}

// Writer allows you to write urls into persistent storage.
type Writer interface {
	// WriteURL to storage, returns short Record.
	WriteURL(url string, uid string) (string, error)
}

type BatchWriter interface {
	BatchWrite(uid string, in []Record) ([]Record, error)
}

type BatchRemover interface {
	BatchRemove(uid string, ids ...string) error
}

type RecordID string

type Record struct {
	ID            string
	ShortURL      string
	OriginalURL   string
	CorrelationID string
}
