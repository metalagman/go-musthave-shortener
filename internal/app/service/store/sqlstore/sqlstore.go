package sqlstore

import (
	"database/sql"
	"fmt"
	"runtime"
	"shortener/internal/app/logger"
	"shortener/pkg/workerpool"
)

type Store struct {
	baseURL string
	base    int
	db      *sql.DB
	log     logger.Logger

	wp *workerpool.Pool
}

// New constructor
func New(db *sql.DB, opts ...Option) (*Store, error) {
	const (
		defaultBase = 36
	)

	s := &Store{
		base: defaultBase,
		db:   db,
		wp:   workerpool.New(),
		log:  logger.Global().Component("Store"),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}

type Option func(*Store)

func WithBaseURL(url string) Option {
	return func(s *Store) {
		s.baseURL = url
	}
}

// Start db connection
func (s *Store) Start() error {

	s.wp.Start(runtime.GOMAXPROCS(0) * 2)

	return nil
}

// Stop store db connection
func (s *Store) Stop() error {
	s.wp.Stop()
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("db close: %w", err)
	}
	return nil
}
