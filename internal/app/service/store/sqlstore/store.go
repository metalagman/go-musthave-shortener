package sqlstore

import (
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	pg "github.com/lib/pq"
	"github.com/russianlagman/go-musthave-shortener/internal/app/service/store"
	"log"
	"strconv"
)

// store.Store interface implementation
var _ store.Store = (*Store)(nil)

func (s *Store) ReadURL(id string) (string, error) {
	const readSQL = `
		SELECT original_url FROM urls WHERE id=$1
`
	rawID, err := s.idToInt64(id)
	if err != nil {
		return "", fmt.Errorf("invalid id %q: %w", id, store.ErrBadInput)
	}

	var url string
	err = s.db.QueryRow(readSQL, rawID).Scan(&url)
	if err != nil {
		return "", fmt.Errorf("read url query: %w", err)
	}

	return url, err
}

func (s *Store) WriteURL(url string, uid string) (string, error) {
	const writeSQL = `
		INSERT INTO urls (uid, original_url)
		VALUES ($1, $2)
		RETURNING id
`

	if err := store.ValidateURL(url); err != nil {
		return "", err
	}

	var rawID int64
	err := s.db.QueryRow(writeSQL, uid, url).Scan(&rawID)
	if err != nil {
		var pgErr *pg.Error
		if errors.As(err, &pgErr) {
			if pgerrcode.IsIntegrityConstraintViolation(string(pgErr.Code)) {
				err := s.db.QueryRow(`SELECT id FROM urls WHERE original_url = $1`, url).Scan(&rawID)
				if err != nil {
					return "", fmt.Errorf("query conflicting id: %w", err)
				}
				return "", &store.ConflictError{
					ExistingURL: s.shortURL(s.idFromInt64(rawID)),
				}
			}
		}

		return "", fmt.Errorf("write url query: %w", err)
	}

	return s.shortURL(s.idFromInt64(rawID)), nil
}

func (s *Store) ReadUserData(uid string) []store.Record {
	const readAllSQL = `
		SELECT id, original_url FROM urls WHERE uid=$1
`

	var result []store.Record

	rows, err := s.db.Query(readAllSQL, uid)
	if err != nil {
		return result
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var (
			rawID       int64
			originalURL string
		)
		if err := rows.Scan(&rawID, &originalURL); err != nil {
			log.Printf("scan error: %v", err)
			break
		}
		id := s.idFromInt64(rawID)
		result = append(result, store.Record{
			ID:          id,
			OriginalURL: originalURL,
			ShortURL:    s.shortURL(id),
		})
	}

	return result
}

// idFromInt64 converts sql id to short string id
func (s *Store) idFromInt64(id int64) string {
	return strconv.FormatInt(id, s.base)
}

// idToInt64 converts short string id to sql id
func (s *Store) idToInt64(id string) (int64, error) {
	rawID, err := strconv.ParseInt(id, s.base, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid id %q: %w", id, store.ErrBadInput)
	}
	return rawID, nil
}

// shortURL returns short url of the id
func (s *Store) shortURL(id string) string {
	return fmt.Sprintf("%s/%s", s.baseURL, id)
}
