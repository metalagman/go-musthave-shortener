package main

import (
	"context"
	"errors"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/russianlagman/go-musthave-shortener/internal/app"
	"github.com/russianlagman/go-musthave-shortener/internal/app/handlers/json"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Config struct {
	ListenAddr string `env:"SERVER_ADDRESS,required"`
	BaseUrl    string `env:"BASE_URL,required"`
}

// Load config from environment and from .env file (if exists)
func (config *Config) Load() error {
	err := godotenv.Load()
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	err = env.Parse(config)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		osCall := <-stop
		log.Printf("system call: %+v", osCall)
		cancel()
	}()

	config := Config{}
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if err := serve(ctx, config); err != nil {
		log.Printf("failed to serve: %+v\n", err)
	}
}

func serve(ctx context.Context, config Config) (err error) {
	shortener := app.NewMemoryShortenerService(config.ListenAddr, config.BaseUrl)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/{id:[0-9a-z]+}", app.ReadHandler(shortener))
	r.Post("/api/shorten", json.WriteHandler(shortener))
	r.Post("/", app.WriteHandler(shortener))
	log.Printf("listening on %s\n", config.ListenAddr)

	srv := &http.Server{
		Addr:    config.ListenAddr,
		Handler: r,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%+v\n", err)
		}
	}()

	log.Printf("server started")

	<-ctx.Done()

	log.Printf("server stopped")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("server shutdown failed: %+v", err)
	}

	log.Printf("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}

	return
}
