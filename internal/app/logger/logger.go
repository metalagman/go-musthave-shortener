package logger

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func init() {
	// setup global logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
}

type Logger struct {
	zerolog.Logger
}
type Component interface {
	// LoggerComponent returns component name used in component loggers
	LoggerComponent() string
}

// New constructor
func New(verbose bool) Logger {
	logLevel := zerolog.InfoLevel
	if verbose {
		logLevel = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(logLevel)

	return Logger{log.Logger}
}

// Global returns current global logger
func Global() Logger {
	return Logger{log.Logger}
}

// Ctx creates context logger
func Ctx(ctx context.Context) Logger {
	logger := zerolog.Ctx(ctx)
	return Logger{Logger: *logger}
}

// For creates logger for specified component or returns current logger
func (l Logger) For(in interface{}) Logger {
	if v, ok := in.(Component); ok {
		return l.Component(v.LoggerComponent())
	}
	return l
}

// Component creates child logger for named component
func (l Logger) Component(name string) Logger {
	return Logger{Logger: l.With().Str("component", name).Logger()}
}
