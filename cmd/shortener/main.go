package main

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/russianlagman/go-musthave-shortener/internal/app/service/store/sqlstore"
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
		log.Fatalf("config load error: %v", err)
	}

	if err := c.Validate(); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			log.Fatalf(
				"invalid value %s for config param %s, expected format: %s",
				err.Value(),
				err.StructField(),
				err.ActualTag(),
			)
		}
	}

	if err := serve(ctx, c); err != nil {
		log.Printf("failed to serve: %+v\n", err)
	}
}

func serve(ctx context.Context, config *Config) (err error) {
	s := sqlstore.NewStore(
		sqlstore.WithBaseURL(config.BaseURL),
		sqlstore.WithListenAddr(config.ListenAddr),
		sqlstore.WithDSN(config.DSN),
	)

	if err := s.Start(); err != nil {
		return fmt.Errorf("store serve failed: %w", err)
	}

	srv := NewServer(config, s)
	go func() {
		log.Printf("listening on %s", config.ListenAddr)
		log.Printf("base url %s", config.BaseURL)
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %+v\n", err)
		}
	}()

	log.Printf("server started")
	<-ctx.Done()
	log.Printf("server stopped")

	if err = s.Stop(); err != nil {
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
