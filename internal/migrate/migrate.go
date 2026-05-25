package migrate

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func Up(db *sql.DB, migrationsDir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	return nil
}
