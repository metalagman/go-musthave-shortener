package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/russianlagman/go-musthave-shortener/internal/app/handlers/basic"
	"github.com/russianlagman/go-musthave-shortener/internal/app/handlers/json"
	"github.com/russianlagman/go-musthave-shortener/internal/app/services/shortener"
	flag "github.com/spf13/pflag"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Config struct {
	ListenAddr           string `env:"SERVER_ADDRESS,required" envDefault:"localhost:8080"`
	BaseURL              string `env:"BASE_URL,required" envDefault:"http://localhost:8080"`
	StorageFilePath      string `env:"FILE_STORAGE_PATH,required" envDefault:"urls.gob"`
	StorageFlushInterval time.Duration
}

func NewConfig() *Config {
	return &Config{StorageFlushInterval: time.Second * 1}
}

// Load config from environment and from .env file (if exists) and from flags
func (config *Config) Load() error {
	err := godotenv.Load()
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf(".env load error: %w", err)
	}
	err = env.Parse(config)
	if err != nil {
		return fmt.Errorf("env parse error: %w", err)
	}

	flag.StringVarP(&config.ListenAddr, "listen-addr", "a", config.ListenAddr, "Server address to listen on")
	flag.StringVarP(&config.BaseURL, "base-url", "b", config.BaseURL, "Base URL for shortened links")
	flag.StringVarP(&config.StorageFilePath, "storage-file-path", "f", config.StorageFilePath, "Storage file path")
	flag.Parse()

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

	c := NewConfig()
	if err := c.Load(); err != nil {
		log.Fatal(err)
	}

	if err := serve(ctx, *c); err != nil {
		log.Printf("failed to serve: %+v\n", err)
	}
}

func serve(ctx context.Context, config Config) (err error) {
	store := shortener.NewMemoryStore(config.ListenAddr, config.BaseURL, config.StorageFilePath)

	log.Printf("reading db from %q", config.StorageFilePath)
	if err = store.ReadDB(); err != nil {
		return fmt.Errorf("error reading db: %w", err)
	}
	log.Printf("done reading db")

	ticker := time.NewTicker(config.StorageFlushInterval)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				log.Print("timer writing db")
				_ = store.WriteDB()
			}
		}
	}()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/{id:[0-9a-z]+}", basic.ReadHandler(store))
	r.Post("/api/shorten", json.WriteHandler(store))
	r.Post("/", basic.WriteHandler(store))
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

	log.Printf("writing db to %q", config.StorageFilePath)
	if err = store.WriteDB(); err != nil {
		log.Printf("error writing db: %v", err)
	}
	log.Printf("done writing db")

	return
}
