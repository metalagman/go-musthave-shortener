package config

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
	"io/fs"
	"net/url"
	"time"
)

type AppConfig struct {
	ListenAddr           string `env:"SERVER_ADDRESS,default=localhost:8080" validate:"required,hostname_port"`
	BaseURL              string `env:"BASE_URL,default=http://localhost:8080" validate:"required,base_url"`
	StorageFilePath      string `env:"FILE_STORAGE_PATH,default=urls.gob" validate:"required"`
	SecretKey            string `env:"SECRET_KEY,default=change_me" validate:"required"`
	DSN                  string `env:"DATABASE_DSN"`
	StorageFlushInterval time.Duration
}

// New constructor
func New() *AppConfig {
	const defaultStorageFlushInterval = time.Second * 5

	return &AppConfig{StorageFlushInterval: defaultStorageFlushInterval}
}

// Load config from environment and from .env file (if exists) and from flags
func (c *AppConfig) Load() error {
	if err := godotenv.Load(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf(".env load error: %w", err)
	}

	if err := envdecode.StrictDecode(c); err != nil {
		return fmt.Errorf("env decode: %w", err)
	}

	pflag.StringVarP(&c.ListenAddr, "listen-addr", "a", c.ListenAddr, "Server address to listen on")
	pflag.StringVarP(&c.BaseURL, "base-url", "b", c.BaseURL, "Base URL for shortened links")
	pflag.StringVarP(&c.StorageFilePath, "storage-file-path", "f", c.StorageFilePath, "Storage file path")
	pflag.StringVarP(&c.DSN, "dsn", "d", c.DSN, "Database connection DSN")
	pflag.Parse()

	return nil
}

func (c *AppConfig) Validate() error {
	validate := validator.New()

	_ = validate.RegisterValidation("base_url", func(fl validator.FieldLevel) bool {
		u, err := url.Parse(fl.Field().String())
		return err == nil && u.Scheme != "" && u.Host != ""
	})

	return validate.Struct(c)
}
