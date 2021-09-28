package securecookie

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"hash"
	"net/http"
)

var (
	_ Encoder = (*SecureCookie)(nil)
	_ Decoder = (*SecureCookie)(nil)
)

type SecureCookie struct {
	hashKey  []byte
	hashFunc func() hash.Hash
	hmacFunc func() hash.Hash
}

func New(secretKey string) *SecureCookie {
	hashFunc := sha256.New
	hashKey := []byte(secretKey)
	return &SecureCookie{
		hashFunc: hashFunc,
		hashKey:  hashKey,
		hmacFunc: func() hash.Hash { return hmac.New(hashFunc, hashKey) },
	}
}

// Decode cookie value if the signature is valid
func (sc *SecureCookie) Decode(cookie *http.Cookie) error {
	e := envelope{}
	if err := e.Decode(cookie.Value); err != nil {
		return fmt.Errorf("envelope decode: %w", err)
	}

	mac, err := generateHmac(sc.hmacFunc(), e.Message)
	if err != nil {
		return fmt.Errorf("generate hmac: %w", err)
	}

	if !hmac.Equal(mac, e.Signature) {
		return ErrDecodeError
	}

	cookie.Value = string(e.Message)

	return nil
}

// Encode cookie value into envelope with signature
func (sc *SecureCookie) Encode(cookie *http.Cookie) error {
	e := envelope{
		Message: []byte(cookie.Value),
	}

	mac, err := generateHmac(sc.hmacFunc(), e.Message)
	if err != nil {
		return fmt.Errorf("generate hmac: %w", err)
	}
	e.Signature = mac

	v, err := e.Encode()
	if err != nil {
		return fmt.Errorf("envelope encode: %w", err)
	}

	cookie.Value = v

	return nil
}

// generate hmac using provided hmac function
func generateHmac(h hash.Hash, msg []byte) ([]byte, error) {
	_, err := h.Write(msg)
	if err != nil {
		return nil, fmt.Errorf("hmac write: %w", err)
	}
	return h.Sum(nil), nil
}
