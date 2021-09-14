package auth

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/russianlagman/go-musthave-shortener/internal/app/service/securecookie"
	"log"
	"net/http"
	"time"
)

const uidCookieName string = "uid"

type ContextKeyUid struct{}

func SecureCookie(secretKey string) func(next http.Handler) http.Handler {
	sc := securecookie.New(secretKey)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("auth middleware")
			var uid string
			var err error

			// reading or regenerating uid
			uid, err = readUid(r, sc)
			if err != nil {
				uid, err = regenerateUid(w, sc)
				if err != nil {
					log.Printf("auth middleware error: %v", err)
					next.ServeHTTP(w, r)
					return
				}
			}

			// inject uid into context
			ctx := context.WithValue(r.Context(), ContextKeyUid{}, uid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// readUid from request cookie and store it into a new context
func readUid(r *http.Request, sc *securecookie.SecureCookie) (string, error) {
	cookie, err := r.Cookie(uidCookieName)
	if err != nil {
		return "", fmt.Errorf("cookie read: %w", err)
	}
	if err := sc.Decode(cookie); err != nil {
		return "", fmt.Errorf("auth decode error: %w", err)
	}
	return cookie.Value, nil
}

// regenerateUid and send it within the cookie
func regenerateUid(w http.ResponseWriter, sc *securecookie.SecureCookie) (string, error) {
	uid := uuid.New().String()

	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := &http.Cookie{
		Name:    uidCookieName,
		Value:   uid,
		Expires: expiration,
	}
	err := sc.Encode(cookie)
	if err != nil {
		return "", fmt.Errorf("encode error: %w", err)
	}

	http.SetCookie(w, cookie)

	return uid, nil
}
