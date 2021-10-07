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
		defer func() {
			_ = stmt.Close()
		}()
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

func (s *Store) BatchRemove(uid string, ids ...string) error {
	const softDeleteQuery = `
		UPDATE urls SET deleted_at = NOW()
		WHERE id=$1 and uid=$2
`
	asyncRemoveJob := func(uid string, id string) Job {
		return func() {
			rawID, err := s.idToInt64(id)
			if err != nil {
				s.log.Error().Err(err).Msg("Uint conversion failure")
				return
			}

			if err := s.execQuery(softDeleteQuery, rawID, uid); err != nil {
				s.log.Error().Err(err).Msg("Exec failure")
			}
		}
	}

	go func() {
		for _, v := range ids {
			s.jobs <- asyncRemoveJob(uid, v)
		}
	}()

	return nil
}
