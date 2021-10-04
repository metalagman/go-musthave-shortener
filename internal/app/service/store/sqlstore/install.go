package sqlstore

import (
	"database/sql"
	"fmt"
)

func (s *Store) createTables() error {
	//_, _ = s.db.Exec(`DROP TABLE IF EXISTS "urls"`)
	const (
		sqlCreateTable = `CREATE TABLE IF NOT EXISTS "urls" (
    id BIGSERIAL primary key,
  	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMPTZ,
	uid UUID,
	original_url TEXT 
);`
		sqlUniqueIndex = `CREATE UNIQUE INDEX urls_unique_original_url_null
ON urls(original_url)
WHERE deleted_at IS NULL`
	)

	if err := execQuery(s.db, sqlCreateTable); err != nil {
		return fmt.Errorf("create table: %w", err)
	}

	if err := execQuery(s.db, sqlUniqueIndex); err != nil {
		return fmt.Errorf("create unuque idx: %w", err)
	}

	return nil
}

func execQuery(db *sql.DB, query string) error {
	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}
