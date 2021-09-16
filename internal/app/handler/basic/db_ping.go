package basic

import (
	"context"
	_ "database/sql"
	"github.com/jackc/pgx"
	"log"
	"net/http"
)

func PingHandler(dsn string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		connConfig, err := pgx.ParseURI(dsn)
		if err != nil {
			log.Printf("dsn %q error: %v", dsn, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		db, err := pgx.Connect(connConfig)
		if err != nil {
			log.Printf("db connection error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := db.Ping(context.Background()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
