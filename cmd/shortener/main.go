package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/russianlagman/go-musthave-shortener/internal/app"
	"log"
	"net/http"
)

func main() {
	addr := "localhost:8080"
	shortener := app.NewMemoryShortenerService(addr)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/{id:[0-9a-z]+}", app.ReadHandler(shortener))
	r.Post("/", app.WriteHandler(shortener))
	log.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
