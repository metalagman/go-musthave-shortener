package sqlstore

import (
	"database/sql"
	"fmt"
	"shortener/internal/app/service/store"
)

// store.BatchWriter interface implementation
var _ store.BatchWriter = (*Store)(nil)
var _ store.BatchRemover = (*Store)(nil)

func (s *Store) BatchWrite(uid string, in []store.Record) ([]store.Record, error) {
	const batchWriteSQL = `
INSERT INTO urls (uid, original_url)
VALUES ($1, $2)
RETURNING id
`
	err := s.inTransaction(func(tx *sql.Tx) error {
		stmt, err := tx.Prepare(batchWriteSQL)
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

func (s *Store) BatchRemove(uid string, in []string) error {
	return nil
}
