package app

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func ReadHandler(svc ShortenerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/")
		u, err := svc.Read(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, u, http.StatusTemporaryRedirect)
	}
}

func WriteHandler(svc ShortenerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		u := string(body)
		if !IsUrl(u) {
			http.Error(w, "bad url", http.StatusBadRequest)
			return
		}
		id, err := svc.Write(u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(201)
		_, _ = w.Write([]byte(id))
	}
}
