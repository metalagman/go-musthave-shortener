package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"net/http/pprof"
	"shortener/internal/app/config"
	"shortener/internal/app/handler/api"
	"shortener/internal/app/handler/basic"
	"shortener/internal/app/logger"
	mw "shortener/internal/app/middleware"
	"shortener/internal/app/service/store/sqlstore"
	"shortener/internal/migrate"
	"time"
)

type App struct {
	config *config.AppConfig
	store  *sqlstore.Store
	log    logger.Logger
}

func (a *App) LoggerComponent() string {
	panic("App")
}

func New(config *config.AppConfig, l logger.Logger) (*App, error) {
	db, err := sql.Open("postgres", config.DSN)
	if err != nil {
		return nil, fmt.Errorf("db open: %w", err)
	}

	if err = migrate.Up(db); err != nil {
		return nil, fmt.Errorf("migrate up: %w", err)
	}

	s, err := sqlstore.New(
		db,
		sqlstore.WithBaseURL(config.BaseURL),
	)
	if err != nil {
		return nil, fmt.Errorf("store init: %w", err)
	}

	a := &App{
		config: config,
		store:  s,
		log:    l,
	}

	return a, nil
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
		a.log.Debug().Msgf("Listening on %s", a.config.ListenAddr)
		a.log.Debug().Msgf("Base URL %s", a.config.BaseURL)

		if a.config.EnableHTTPS {
			manager := &autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(a.config.ListenAddr),
			}
			srv.TLSConfig = manager.TLSConfig()
			if err := srv.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
				a.log.Fatal().Err(err).Msg("Socket listen failure")
			}
		} else {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				a.log.Fatal().Err(err).Msg("Socket listen failure")
			}
		}
	}()

	a.log.Debug().Msgf("Server started")
	<-ctx.Done()
	a.log.Debug().Msgf("Server stopped")

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

	a.log.Debug().Msgf("Server exited properly")

	return nil
}

func (a *App) router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(mw.Log(a.log))
	r.Use(mw.SecureCookieAuth(a.config.SecretKey))
	r.Use(mw.GzipResponseWriter)
	r.Use(mw.GzipRequestReader)

	AttachProfiler(r)

	r.With(mw.ContentTypeJSON).Get("/api/user/urls", api.UserDataHandler(a.store))
	r.With(mw.ContentTypeJSON).Post("/api/shorten", api.WriteHandler(a.store))
	r.With(mw.ContentTypeJSON).Post("/api/shorten/batch", api.BatchWriteHandler(a.store))
	r.With(mw.ContentTypeJSON).Delete("/api/user/urls", api.BatchRemoveHandler(a.store))
	r.With(mw.ContentTypeJSON, mw.TrustedNetwork(a.config.TrustedNetwork)).Get("/api/internal/stats", api.StatHandler(a.store))
	r.Get("/{id:[0-9a-z]+}", basic.ReadHandler(a.store))
	r.Post("/", basic.WriteHandler(a.store))
	r.Get("/ping", basic.PingHandler(a.store))

	return r
}

func AttachProfiler(router *chi.Mux) {
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	// Manually add support for paths linked to by index page at /debug/pprof/
	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/debug/pprof/block", pprof.Handler("block"))
}
