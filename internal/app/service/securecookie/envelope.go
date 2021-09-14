package securecookie

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
)

type envelope struct {
	Message   []byte
	Signature []byte
}

// Encode envelope into base64-hmac gob string
func (e *envelope) Encode() (string, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(e); err != nil {
		return "", fmt.Errorf("gob encoder: %w", err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// Decode envelope from base64-hmac gob string
func (e *envelope) Decode(in string) error {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return fmt.Errorf("base64 decode: %w", err)
	}
	dec := gob.NewDecoder(bytes.NewReader(b))
	if err := dec.Decode(e); err != nil {
		return fmt.Errorf("gob decode: %w", err)
	}
	return nil
}
