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
