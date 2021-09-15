package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/russianlagman/go-musthave-shortener/internal/app/handler/basic"
	"github.com/russianlagman/go-musthave-shortener/internal/app/handler/json"
	app "github.com/russianlagman/go-musthave-shortener/internal/app/middleware"
	"github.com/russianlagman/go-musthave-shortener/internal/app/service/store"
	"net/http"
)

func NewServer(config *Config, store store.Store) *http.Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(app.SecureCookieAuth("test secret"))
	r.Use(app.GzipResponseWriter)
	r.Use(app.GzipRequestReader)
	r.With(app.JsonNegotiator).Get("/user/urls", json.UserDataHandler(store))
	r.With(app.JsonNegotiator).Post("/api/shorten", json.WriteHandler(store))
	r.Get("/{id:[0-9a-z]+}", basic.ReadHandler(store))
	r.Post("/", basic.WriteHandler(store))

	return &http.Server{
		Addr:    config.ListenAddr,
		Handler: r,
	}
}
