package sqlstore

import (
	"fmt"
	"github.com/russianlagman/go-musthave-shortener/internal/app/service/store"
)

// store.HealthChecker interface implementation
var _ store.HealthChecker = (*Store)(nil)

func (s *Store) HealthCheck() error {
	if err := s.db.Ping(); err != nil {
		return fmt.Errorf("ping: %w", err)
	}
	return nil
}
