package sqlstore

import (
	"database/sql"
	"fmt"
)

type Store struct {
	baseURL string
	dsn     string
	base    int

	db *sql.DB
}

// New constructor
func New(opts ...Option) *Store {
	const (
		defaultBase = 36
	)

	s := &Store{
		base: defaultBase,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

type Option func(*Store)

func WithBaseURL(url string) Option {
	return func(s *Store) {
		s.baseURL = url
	}
}

func WithDSN(dsn string) Option {
	return func(s *Store) {
		s.dsn = dsn
	}
}

// Start db connection
func (s *Store) Start() error {
	var err error
	s.db, err = sql.Open("postgres", s.dsn)
	if err != nil {
		return fmt.Errorf("db connection: %w", err)
	}
	if err := s.createTables(); err != nil {
		return fmt.Errorf("create tables: %w", err)
	}
	return nil
}

// Stop store db connection
func (s *Store) Stop() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("db close: %w", err)
	}
	return nil
}
