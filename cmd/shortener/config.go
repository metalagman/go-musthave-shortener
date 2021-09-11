package main

import (
	"errors"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
	"io/fs"
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

	pflag.StringVarP(&config.ListenAddr, "listen-addr", "a", config.ListenAddr, "Server address to listen on")
	pflag.StringVarP(&config.BaseURL, "base-url", "b", config.BaseURL, "Base URL for shortened links")
	pflag.StringVarP(&config.StorageFilePath, "storage-file-path", "f", config.StorageFilePath, "Storage file path")
	pflag.Parse()

	return nil
}
