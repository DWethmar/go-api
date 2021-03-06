package database

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	// Import from file system.
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations runs the migration from a folder.
func RunMigrations(db *sql.DB, dbName, folder string, version int) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+folder, // "file:///migrations",
		dbName,           // "postgres",
		driver,
	)

	if err != nil {
		return err
	}

	v, dirty, _ := m.Version()

	fmt.Printf("migration current version: %v, dirty: %v.", v, dirty)

	return m.Up()
}
