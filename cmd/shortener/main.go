package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/russianlagman/go-musthave-shortener/internal/app"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	log.Printf("server started")

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		osCall := <-stop
		log.Printf("system call: %+v", osCall)
		cancel()
	}()

	if err := serve(ctx); err != nil {
		log.Printf("failed to serve: %+v\n", err)
	}
}

func serve(ctx context.Context) (err error) {
	addr := "localhost:8080"
	shortener := app.NewMemoryShortenerService(addr)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/{id:[0-9a-z]+}", app.ReadHandler(shortener))
	r.Post("/", app.WriteHandler(shortener))
	log.Printf("listening on %s\n", addr)

	srv := &http.Server{
		Addr:    addr,
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
