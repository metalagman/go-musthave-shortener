package main

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"os/signal"
	"shortener/internal/app"
	"shortener/internal/app/config"
	"shortener/internal/app/logger"
	"shortener/pkg/version"
	"syscall"
)

func main() {
	fmt.Println(version.Print())

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

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
