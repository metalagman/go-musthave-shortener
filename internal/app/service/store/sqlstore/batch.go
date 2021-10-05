package sqlstore

import (
	"database/sql"
	"fmt"
	"log"
	"shortener/internal/app/service/store"
	"strconv"
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
	type RemoveRequest struct {
		uid string
		id  string
	}

	q := make(chan RemoveRequest)

	go func() {
		for _, v := range ids {
			q <- RemoveRequest{
				uid,
				v,
			}
		}
		close(q)
	}()

	const softDeleteQuery = `
		UPDATE urls SET deleted_at = NOW()
		WHERE id=$1 and uid=$2
`
	for i := 0; i < s.workerNum; i++ {
		go func(id int, in <-chan RemoveRequest) {
			log.Printf("worker [%d] started", id)
			for v := range in {
				id, err := strconv.ParseUint(v.id, s.base, 64)
				if err != nil {
					log.Printf("parse uint: %v", err)
					continue
				}

				if err := s.execQuery(softDeleteQuery, v.id, id); err != nil {
					log.Printf("exec: %v", err)
				}
			}
			log.Printf("worker [%d] finished", id)
		}(i, q)
	}

	return nil
}
