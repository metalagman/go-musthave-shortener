package basic

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func PingHandler(URI string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("postgres", URI)
		if err != nil {
			log.Printf("db connection error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		defer func() {
			_ = db.Close()
		}()

		if err := db.Ping(); err != nil {
			log.Printf("ping error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
