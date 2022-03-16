package basic

import (
	_ "github.com/lib/pq"
	"net/http"
	"shortener/internal/app/logger"
	"shortener/internal/app/service/store"
)

// PingHandler allows you to check system health.
//
//	curl -v http://localhost:8080/ping
func PingHandler(s store.HealthChecker) http.HandlerFunc {
	log := logger.Global().Component("Handler::Ping")
	return func(w http.ResponseWriter, r *http.Request) {
		if err := s.HealthCheck(); err != nil {
			log.Error().Err(err).Msg("DB ping failure")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
