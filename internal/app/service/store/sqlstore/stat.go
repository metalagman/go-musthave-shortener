package sqlstore

import (
	"shortener/internal/app/service/store"
)

// store.StatProvider interface implementation
var _ store.StatProvider = (*Store)(nil)

func (s *Store) Stat() (store.StatData, error) {
	res := store.StatData{
		URLCount:  0,
		UserCount: 0,
	}

	return res, nil
}
