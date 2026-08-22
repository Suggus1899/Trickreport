package db

import (
	"embed"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// RunMigrations applies all pending migrations to the database.
// databaseURL must be a valid pgx connection string.
//
// The pgx/v5 driver used by golang-migrate expects the URL scheme to be
// "pgx5://". If the supplied DSN uses the more common "postgres://" scheme
// (as pgxpool and the rest of the application do), it is rewritten here so
// callers can pass the same DATABASE_URL everywhere.
func RunMigrations(databaseURL string) error {
	d, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}

	migrateURL := strings.Replace(databaseURL, "postgres://", "pgx5://", 1)

	m, err := migrate.NewWithSourceInstance("iofs", d, migrateURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

// RollbackLastMigration rolls back the most recently applied migration.
// databaseURL must be a valid pgx connection string (see RunMigrations).
func RollbackLastMigration(databaseURL string) error {
	d, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}

	migrateURL := strings.Replace(databaseURL, "postgres://", "pgx5://", 1)

	m, err := migrate.NewWithSourceInstance("iofs", d, migrateURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to roll back migration: %w", err)
	}
	return nil
}
