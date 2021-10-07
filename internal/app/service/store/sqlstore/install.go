package sqlstore

import (
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
		sqlUniqueIndex = `CREATE UNIQUE INDEX IF NOT EXISTS urls_unique_original_url_null
ON urls(original_url)
WHERE deleted_at IS NULL`
	)

	if err := s.execQuery(sqlCreateTable); err != nil {
		return fmt.Errorf("create table: %w", err)
	}

	if err := s.execQuery(sqlUniqueIndex); err != nil {
		return fmt.Errorf("create unuque idx: %w", err)
	}

	return nil
}
