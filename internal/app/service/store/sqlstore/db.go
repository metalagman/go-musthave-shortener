package sqlstore

import (
	"database/sql"
	"fmt"
)

func (s *Store) inTransaction(cb func(tx *sql.Tx) error) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("tx begin: %w", err)
	}

	if err := cb(tx); err != nil {
		err = fmt.Errorf("callback: %w", err)
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("rollback: %w", err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("tx commit: %w", err)
	}

	return nil
}

func (s *Store) execQuery(query string, args ...interface{}) error {
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.Exec(args...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}
