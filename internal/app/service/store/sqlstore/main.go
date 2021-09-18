package sqlstore

import (
	"database/sql"
	"fmt"
)

type Store struct {
	listenAddr string
	baseURL    string
	dsn        string
	base       int

	db *sql.DB
}

// NewStore constructor
func NewStore(opts ...StoreOption) *Store {
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

type StoreOption func(*Store)

func WithListenAddr(addr string) StoreOption {
	return func(s *Store) {
		s.listenAddr = addr
	}
}

func WithBaseURL(url string) StoreOption {
	return func(s *Store) {
		s.baseURL = url
	}
}

func WithDSN(dsn string) StoreOption {
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

func (s *Store) createTables() error {
	sql := `CREATE TABLE IF NOT EXISTS "urls" (
    id BIGSERIAL primary key,
  	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	uid UUID,
	original_url TEXT
);`
	stmt, err := s.db.Prepare(sql)
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}
