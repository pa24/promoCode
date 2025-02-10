package migrations

import (
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
)

func RunMigrations(db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("dialect error: %v", err)
	}

	migrationsDir := "./internal/storage/migrations"
	if err := goose.Up(db, migrationsDir); err != nil {
		return err
	}

	return nil
}
