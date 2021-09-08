package main

import (
	"context"
	"encoding/gob"
	"errors"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/russianlagman/go-musthave-shortener/internal/app/handlers/basic"
	"github.com/russianlagman/go-musthave-shortener/internal/app/handlers/json"
	"github.com/russianlagman/go-musthave-shortener/internal/app/services/shortener"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Config struct {
	ListenAddr      string `env:"SERVER_ADDRESS,required" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL,required" envDefault:"http://localhost:8080"`
	StorageFilePath string `env:"FILE_STORAGE_PATH,required" envDefault:"urls.gob"`
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

	c := Config{}
	err := c.Load()
	if err != nil {
		log.Fatal(err)
	}

	if err := serve(ctx, c); err != nil {
		log.Printf("failed to serve: %+v\n", err)
	}
}

func serve(ctx context.Context, config Config) (err error) {
	store := shortener.NewMemoryStore(config.ListenAddr, config.BaseURL)
	store.SetDB(readDb(config.StorageFilePath))
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

	return
}

// readDb from file at filePath
func readDb(filePath string) shortener.MemoryDB {
	db := make(shortener.MemoryDB)

	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatalf("error reading db at %q: %v", filePath, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	decoder := gob.NewDecoder(file)

	err = decoder.Decode(&db)
	if err != nil && err != io.EOF {
		log.Fatalf("decode error: %v", err)
	}

	return db
}
