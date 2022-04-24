package sqlstore

import (
	"fmt"
	"shortener/internal/app/service/store"
)

// store.StatProvider interface implementation
var _ store.StatProvider = (*Store)(nil)

func (s *Store) Stat() (*store.StatData, error) {
	const (
		userSQL = `SELECT COUNT(distinct uid) FROM urls`
		urlSQL  = `SELECT COUNT(*) FROM urls`
	)

	res := &store.StatData{}

	if err := s.db.QueryRow(userSQL).Scan(&res.UserCount); err != nil {
		return nil, fmt.Errorf("user count query: %w", err)
	}

	if err := s.db.QueryRow(urlSQL).Scan(&res.URLCount); err != nil {
		return nil, fmt.Errorf("url count query: %w", err)
	}

	return res, nil
}
