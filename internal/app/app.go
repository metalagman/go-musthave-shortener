package app

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"shortener/internal/app/config"
	"shortener/internal/app/handler/basic"
	"shortener/internal/app/handler/json"
	mw "shortener/internal/app/middleware"
	"shortener/internal/app/service/store/sqlstore"
	"time"
)

type App struct {
	config *config.AppConfig
	store  *sqlstore.Store
}

func New(config *config.AppConfig) *App {
	return &App{
		config: config,
		store: sqlstore.New(
			sqlstore.WithBaseURL(config.BaseURL),
			sqlstore.WithDSN(config.DSN),
		),
	}
}

func (a *App) Serve(ctx context.Context) error {
	if err := a.store.Start(); err != nil {
		return fmt.Errorf("store start: %w", err)
	}

	srv := &http.Server{
		Addr:    a.config.ListenAddr,
		Handler: a.router(),
	}

	go func() {
		log.Printf("listening on %s", a.config.ListenAddr)
		log.Printf("base url %s", a.config.BaseURL)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %+v\n", err)
		}
	}()

	log.Printf("server started")
	<-ctx.Done()
	log.Printf("server stopped")

	if err := a.store.Stop(); err != nil {
		return fmt.Errorf("store shutdown: %w", err)
	}

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	log.Printf("server exited properly")

	return nil
}

func (a *App) router() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(mw.SecureCookieAuth(a.config.SecretKey))
	r.Use(mw.GzipResponseWriter)
	r.Use(mw.GzipRequestReader)

	r.Get("/user/urls", json.UserDataHandler(a.store))
	r.With(mw.ContentTypeJSON).Post("/api/shorten", json.WriteHandler(a.store))
	r.With(mw.ContentTypeJSON).Post("/api/shorten/batch", json.BatchWriteHandler(a.store))
	r.With(mw.ContentTypeJSON).Delete("/api/user/urls", json.BatchRemoveHandler(a.store))
	r.Get("/{id:[0-9a-z]+}", basic.ReadHandler(a.store))
	r.Post("/", basic.WriteHandler(a.store))
	r.Get("/ping", basic.PingHandler(a.store))

	return r
}
