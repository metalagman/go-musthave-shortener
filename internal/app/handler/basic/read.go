package basic

import (
	"errors"
	"net/http"
	"shortener/internal/app/service/store"
	"strings"
)

// ReadHandler allows you to read short url.
//
//	curl -v http://localhost:8080/xxx
func ReadHandler(s store.Reader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/")
		u, err := s.ReadURL(id)
		if err != nil {
			if errors.Is(err, store.ErrDeleted) {
				http.Error(w, err.Error(), http.StatusGone)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, u, http.StatusTemporaryRedirect)
	}
}
