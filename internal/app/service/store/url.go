package store

import (
	"fmt"
	"net/url"
)

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
