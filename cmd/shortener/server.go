package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/russianlagman/go-musthave-shortener/internal/app/handler/basic"
	"github.com/russianlagman/go-musthave-shortener/internal/app/handler/json"
	app "github.com/russianlagman/go-musthave-shortener/internal/app/middleware"
	"github.com/russianlagman/go-musthave-shortener/internal/app/service/store/memorystore"
	"net/http"
)

func NewServer(config *Config, store *memorystore.Store) *http.Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(app.SecureCookieAuth("test secret"))
	r.Use(app.GzipResponseWriter)
	r.Use(app.GzipRequestReader)
	r.Get("/user/urls", json.UserDataHandler(store))
	r.With(app.ContentTypeJSON).Post("/api/shorten", json.WriteHandler(store))
	r.Get("/{id:[0-9a-z]+}", basic.ReadHandler(store))
	r.Post("/", basic.WriteHandler(store))
	r.Get("/ping", basic.PingHandler(config.DSN))

	return &http.Server{
		Addr:    config.ListenAddr,
		Handler: r,
	}
}
