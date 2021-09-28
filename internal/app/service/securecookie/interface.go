package securecookie

import (
	"errors"
	"net/http"
)

var ErrDecodeError = errors.New("decode error")

// Encoder encodes cookie value
type Encoder interface {
	// Encode cookie
	Encode(cookie *http.Cookie) error
}

// Decoder validates cookie value and returns ErrDecodeError if provided cookie is invalid
type Decoder interface {
	// Decode signed cookie
	Decode(cookie *http.Cookie) error
}
