package main

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func gzipRequestReader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент поддерживает gzip-сжатие
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// создаём gzip.Reader поверх текущего r.Body
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			log.Printf("error gzip reading request body: %v", err)
			return
		}
		defer func(gz *gzip.Reader) {
			_ = gz.Close()
		}(gz)

		// replace body with gzip reader
		r.Body = ioutil.NopCloser(gz)

		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(w, r)
	})
}

func gzipResponseWriter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент поддерживает gzip-сжатие
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}
		defer func(gz *gzip.Writer) {
			_ = gz.Close()
		}(gz)

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}
