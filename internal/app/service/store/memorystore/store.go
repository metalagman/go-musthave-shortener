package memorystore

import (
	"fmt"
	"shortener/internal/app/service/store"
	"strconv"
)

// store.Store interface implementation
var _ store.Store = (*Store)(nil)

func (s *Store) ReadURL(id string) (string, error) {
	intID, err := strconv.ParseUint(id, s.base, 64)
	if err != nil {
		return "", fmt.Errorf("invalid id %q: %w", id, store.ErrBadInput)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	if val, ok := s.db[intID]; ok {
		return val.OriginalURL, nil
	}

	return "", store.ErrNotFound
}

func (s *Store) WriteURL(url string, uid string) (string, error) {
	if err := store.ValidateURL(url); err != nil {
		return "", err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter++
	id := strconv.FormatUint(s.counter, s.base)
	shortURL := fmt.Sprintf("%s/%s", s.baseURL, id)

	s.db[s.counter] = dbRow{
		ID:          id,
		OriginalURL: url,
		ShortURL:    shortURL,
		UID:         uid,
	}

	return shortURL, nil
}

func (s *Store) ReadUserData(uid string) []store.Record {
	var result []store.Record
	for _, row := range s.db {
		if row.UID != uid {
			continue
		}
		result = append(result, store.Record{
			OriginalURL: row.OriginalURL,
			ShortURL:    row.ShortURL,
		})
	}
	return result
}
