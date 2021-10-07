package main

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"shortener/internal/app"
	"shortener/internal/app/config"
	"shortener/internal/app/logger"
)

func main() {
	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		osCall := <-stop
		log.Debug().Msgf("System call: %+v", osCall)
		cancel()
	}()

	c := config.New()
	if err := c.Load(); err != nil {
		log.Fatal().Err(err).Msg("Config load failure")
	}

	l := logger.New(c.Verbose)

	if err := c.Validate(); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			l.Fatal().Err(err).Msgf(
				"Invalid value %s for config param %s, expected format: %s",
				err.Value(),
				err.StructField(),
				err.ActualTag(),
			)
		}
	}

	a, err := app.New(c, l)
	if err != nil {
		log.Fatal().Err(err).Msg("Application init failure")
	}

	if err := a.Serve(ctx); err != nil {
		log.Fatal().Err(err).Msg("Application failed to serve")
	}
}
