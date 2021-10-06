package middleware

import (
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"shortener/internal/app/logger"
	"strings"
)

func GzipRequestReader(next http.Handler) http.Handler {
	log := logger.Global().Component("Middleware::GzipRequestReader")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			log.Error().Err(err).Msg("GZIP reading request body failure")
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
