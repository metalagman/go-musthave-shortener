package store

import (
	"errors"
	"fmt"
	"net/url"
)

var ErrBadInput = errors.New("bad input")
var ErrEmptyInput = fmt.Errorf("%w: empty url", ErrBadInput)
var ErrNotFound = errors.New("not found")

// Store of the db served by store
type Store interface {
	// WriteURL to storage, returns short URL
	WriteURL(url string) (string, error)
	// ReadURL from storage
	ReadURL(id string) (string, error)
}

// ValidateURL checks if input is valid url
func ValidateURL(str string) error {
	if str == "" {
		return ErrEmptyInput
	}
	u, err := url.Parse(str)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("%w: invalid url %q", ErrBadInput, str)
	}
	return nil
}
