package basic

import (
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"shortener/internal/app/service/store"
)

func PingHandler(s store.HealthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := s.HealthCheck(); err != nil {
			log.Printf("ping handler: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
