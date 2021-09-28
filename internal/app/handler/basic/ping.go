package basic

import (
	_ "github.com/lib/pq"
	"github.com/russianlagman/go-musthave-shortener/internal/app/service/store"
	"log"
	"net/http"
)

func PingHandler(s store.Ping) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := s.Ping(); err != nil {
			log.Printf("ping handler: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
