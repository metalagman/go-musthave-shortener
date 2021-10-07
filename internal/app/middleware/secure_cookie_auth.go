package middleware

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"shortener/internal/app/handler"
	"shortener/internal/app/logger"
	"shortener/internal/app/service/securecookie"
	"time"
)

const cookieNameUID string = "uid"

func SecureCookieAuth(secretKey string) func(next http.Handler) http.Handler {
	log := logger.Global().Component("Middleware::SecureCookieAuth")
	sc := securecookie.New(secretKey)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var uid string
			var err error

			// reading or regenerating uid
			uid, err = readUID(r, sc)
			if err != nil {
				if err != http.ErrNoCookie {
					log.Error().Err(err).Msg("UID read failure")
				}
				uid, err = regenerateUID(w, sc)
				if err != nil {
					log.Error().Err(err).Msg("UID regenerate failure")
					next.ServeHTTP(w, r)
					return
				}
			}

			// inject uid into context
			ctx := context.WithValue(r.Context(), handler.ContextKeyUID{}, uid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// readUID from request cookie and store it into a new context
func readUID(r *http.Request, sc *securecookie.SecureCookie) (string, error) {
	cookie, err := r.Cookie(cookieNameUID)
	if err != nil {
		return "", fmt.Errorf("cookie read: %w", err)
	}
	if err := sc.Decode(cookie); err != nil {
		return "", fmt.Errorf("auth decode error: %w", err)
	}
	return cookie.Value, nil
}

// regenerateUID and send it within the cookie
func regenerateUID(w http.ResponseWriter, sc *securecookie.SecureCookie) (string, error) {
	uid := uuid.New().String()

	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := &http.Cookie{
		Name:    cookieNameUID,
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
