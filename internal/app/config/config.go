/*
Package config provides application config structure and tools.
*/
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
	"io/fs"
	"net/url"
	"os"
	"time"
)

/**
{
    "server_address": "localhost:80", // аналог переменной окружения SERVER_ADDRESS или флага -a
    "base_url": "http://localhost", // аналог переменной окружения BASE_URL или флага -b
    "file_storage_path": "/path/to/file.db", // аналог переменной окружения FILE_STORAGE_PATH или флага -f
    "database_dsn": "", // аналог переменной окружения DATABASE_DSN или флага -d
    "enable_https": true // аналог переменной окружения ENABLE_HTTPS или флага -s
}
*/
type AppConfig struct {
	ListenAddr           string `env:"SERVER_ADDRESS,default=localhost:8080" validate:"required,hostname_port" json:"server_address"`
	BaseURL              string `env:"BASE_URL,default=http://localhost:8080" validate:"required,base_url" json:"base_url"`
	StorageFilePath      string `env:"FILE_STORAGE_PATH,default=urls.gob" validate:"required" json:"file_storage_path"`
	SecretKey            string `env:"SECRET_KEY,default=change_me" validate:"required"`
	DSN                  string `env:"DATABASE_DSN" json:"database_dsn"`
	StorageFlushInterval time.Duration
	Verbose              bool   `env:"APP_VERBOSE,default=0"`
	EnableHTTPS          bool   `env:"ENABLE_HTTPS,default=0" json:"enable_https"`
	ConfigFile           string `env:"CONFIG"`
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

	// load config
	if err := envdecode.StrictDecode(c); err != nil {
		return fmt.Errorf("env decode: %w", err)
	}

	// parse config related flags
	pfc := pflag.NewFlagSet("config", pflag.ContinueOnError)
	pfc.StringVarP(&c.ConfigFile, "config", "c", c.ConfigFile, "Config file")
	_ = pfc.Parse(os.Args[1:])
	pflag.CommandLine.AddFlagSet(pfc)

	if c.ConfigFile != "" {
		b, err := os.ReadFile(c.ConfigFile)
		if err != nil {
			return fmt.Errorf("config file read: %w", err)
		}
		if err = json.Unmarshal(b, c); err != nil {
			return fmt.Errorf("config file parse: %w", err)
		}
	}

	// load env one more time
	if err := envdecode.StrictDecode(c); err != nil {
		return fmt.Errorf("env decode: %w", err)
	}

	pflag.StringVarP(&c.ListenAddr, "listen-addr", "a", c.ListenAddr, "Server address to listen on")
	pflag.StringVarP(&c.BaseURL, "base-url", "b", c.BaseURL, "Base URL for shortened links")
	pflag.StringVarP(&c.StorageFilePath, "storage-file-path", "f", c.StorageFilePath, "Storage file path")
	pflag.StringVarP(&c.DSN, "dsn", "d", c.DSN, "Database connection DSN")
	pflag.BoolVarP(&c.Verbose, "verbose", "v", c.Verbose, "Verbose output")
	pflag.BoolVarP(&c.EnableHTTPS, "secure", "s", c.EnableHTTPS, "Enable HTTPS")
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
