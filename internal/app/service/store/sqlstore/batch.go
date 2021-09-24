package sqlstore

import (
	"database/sql"
	"fmt"
	"github.com/russianlagman/go-musthave-shortener/internal/app/service/store"
)

// store.Store interface implementation
var _ store.BatchWriter = (*Store)(nil)

const _batchWriteSQL = `
INSERT INTO urls (uid, original_url)
VALUES ($1, $2)
`

func (s *Store) BatchWrite(uid string, in []store.Record) ([]store.Record, error) {
	err := s.inTransaction(func(tx *sql.Tx) error {
		stmt, err := tx.Prepare(_batchWriteSQL)
		if err != nil {
			return fmt.Errorf("sql prepare: %w", err)
		}
		for i := range in {
			res, err := stmt.Exec(uid, in[i].OriginalURL)
			if err != nil {
				return fmt.Errorf("exec: %w", err)
			}
			id, err := res.LastInsertId()
			if err != nil {
				return fmt.Errorf("last insert id: %w", err)
			}
			in[i].ID = s.idFromInt64(id)
			in[i].ShortURL = s.shortUrl(in[i].ID)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("tx: %w", err)
	}

	return in, nil
}
