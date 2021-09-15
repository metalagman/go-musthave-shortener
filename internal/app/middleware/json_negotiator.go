package middleware

import (
	"net/http"
)

func JsonNegotiator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "invalid content type", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
		if w.Header().Get("ContentType") != "" {
			r.Header.Set("Content-Type", "application/json")
		}
	})
}
