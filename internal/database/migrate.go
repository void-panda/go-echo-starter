package database

import (
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"

	"go-echo-starter/pkg/logger"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Migrator handles database migrations
type Migrator struct {
	m   *migrate.Migrate
	log *logger.Logger
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *sqlx.DB, log *logger.Logger) (*Migrator, error) {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to create migration source: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	return &Migrator{m: m, log: log}, nil
}

// Up runs all pending migrations
func (m *Migrator) Up() error {
	m.log.Info().Msg("Running database migrations...")

	if err := m.m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	version, dirty, err := m.m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	m.log.Info().Uint("version", version).Bool("dirty", dirty).Msg("Migration completed")
	return nil
}

// Down rolls back all migrations
func (m *Migrator) Down() error {
	m.log.Info().Msg("Rolling back all migrations...")

	if err := m.m.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("rollback failed: %w", err)
	}

	m.log.Info().Msg("Rollback completed")
	return nil
}

// Steps migrates up or down by n steps
func (m *Migrator) Steps(n int) error {
	if err := m.m.Steps(n); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration steps failed: %w", err)
	}
	return nil
}

// Version returns current migration version
func (m *Migrator) Version() (uint, bool, error) {
	return m.m.Version()
}

// Close closes the migrator
func (m *Migrator) Close() error {
	sourceErr, dbErr := m.m.Close()
	if sourceErr != nil {
		return sourceErr
	}
	return dbErr
}
