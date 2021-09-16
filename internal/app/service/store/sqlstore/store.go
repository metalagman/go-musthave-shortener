package sqlstore

import "github.com/russianlagman/go-musthave-shortener/internal/app/service/store"

// store.Store interface implementation
var _ store.Store = (*Store)(nil)

func (s *Store) ReadURL(id string) (string, error) {
	return "", nil
}

func (s *Store) WriteURL(url string, uid string) (string, error) {
	return "", nil
}

func (s *Store) ReadUserURLs(uid string) []store.StoredURL {
	var result []store.StoredURL
	return result
}
