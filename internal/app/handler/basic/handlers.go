package basic

import (
	"fmt"
	"github.com/russianlagman/go-musthave-shortener/internal/app/middleware"
	"github.com/russianlagman/go-musthave-shortener/internal/app/service/store"
	"io/ioutil"
	"net/http"
	"strings"
)

func ReadHandler(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/")
		u, err := store.ReadURL(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, u, http.StatusTemporaryRedirect)
	}
}

func WriteHandler(store store.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("error reading body: %v", err), http.StatusBadRequest)
			return
		}
		u := string(body)
		redirectURL, err := store.WriteUserURL(u, r.Context().Value(middleware.ContextKeyUID{}).(string))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(redirectURL))
	}
}
