package basic

import (
	"github.com/russianlagman/go-musthave-shortener/internal/app/service/store"
	"net/http"
	"strings"
)

func ReadHandler(s store.Reader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/")
		u, err := s.ReadURL(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, u, http.StatusTemporaryRedirect)
	}
}
