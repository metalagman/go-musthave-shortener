package main

import (
	"context"
	"github.com/go-playground/validator/v10"
	"log"
	"os"
	"os/signal"
	"shortener/internal/app"
	"shortener/internal/app/config"
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

	c := config.New()
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

	a := app.New(c)

	if err := a.Serve(ctx); err != nil {
		log.Printf("failed to serve: %+v\n", err)
	}
}
