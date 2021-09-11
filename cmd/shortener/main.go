package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/russianlagman/go-musthave-shortener/internal/app/handlers/basic"
	"github.com/russianlagman/go-musthave-shortener/internal/app/handlers/json"
	"github.com/russianlagman/go-musthave-shortener/internal/app/services/shortener"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

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
	store := shortener.NewMemoryStore(
		config.ListenAddr,
		config.BaseURL,
		config.StorageFilePath,
		config.StorageFlushInterval,
	)

	if err := store.Serve(); err != nil {
		return fmt.Errorf("store serve failed: %w", err)
	}

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
			log.Fatalf("listen: %+v\n", err)
		}
	}()

	log.Printf("server started")

	<-ctx.Done()

	log.Printf("server stopped")

	if err = store.Shutdown(); err != nil {
		return fmt.Errorf("store shutdown failed: %w", err)
	}

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctxShutdown); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	log.Printf("server exited properly")

	return
}
