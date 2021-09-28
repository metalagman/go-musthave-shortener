package middleware

import (
	"compress/gzip"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func GzipRequestReader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			log.Printf("error gzip reading request body: %v", err)
			return
		}
		defer func() {
			_ = gz.Close()
			_ = r.Body.Close()
		}()

		// replace body with gzip reader
		r.Body = ioutil.NopCloser(gz)

		next.ServeHTTP(w, r)
	})
}
