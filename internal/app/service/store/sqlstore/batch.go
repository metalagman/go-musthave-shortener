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
RETURNING id`

func (s *Store) BatchWrite(uid string, in []store.Record) ([]store.Record, error) {
	err := s.inTransaction(func(tx *sql.Tx) error {
		stmt, err := tx.Prepare(_batchWriteSQL)
		if err != nil {
			return fmt.Errorf("sql prepare: %w", err)
		}
		defer func(stmt *sql.Stmt) {
			_ = stmt.Close()
		}(stmt)
		for i := range in {
			var rawID int64
			err := stmt.QueryRow(uid, in[i].OriginalURL).Scan(&rawID)
			if err != nil {
				return fmt.Errorf("sql query: %w", err)
			}
			in[i].ID = s.idFromInt64(rawID)
			in[i].ShortURL = s.shortURL(in[i].ID)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("tx: %w", err)
	}

	return in, nil
}
