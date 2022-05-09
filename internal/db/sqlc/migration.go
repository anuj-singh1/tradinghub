package db

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	MIGRATE_DB             = "tradinghub"
	MIGRATE_FILES          = "file://./migration"
	NO_CHANGE_IN_MIGRATION = "no change"
)

// migration steps.
func RunMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		MIGRATE_FILES, MIGRATE_DB, driver)
	if err != nil {
		return err
	}

	err = m.Up()

	if err != nil && err.Error() != NO_CHANGE_IN_MIGRATION {
		return err
	}

	return nil
}
