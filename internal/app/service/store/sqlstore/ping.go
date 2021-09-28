package sqlstore

import (
	"fmt"
	"github.com/russianlagman/go-musthave-shortener/internal/app/service/store"
)

// store.Ping interface implementation
var _ store.Ping = (*Store)(nil)

func (s *Store) Ping() error {
	if err := s.db.Ping(); err != nil {
		return fmt.Errorf("ping: %w", err)
	}
	return nil
}
