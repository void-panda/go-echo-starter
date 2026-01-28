package database

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"go-echo-starter/internal/config"
	"go-echo-starter/pkg/logger"
)

// PostgreSQL holds the database connection
type PostgreSQL struct {
	DB *sqlx.DB
}

// NewPostgreSQL creates a new PostgreSQL connection
func NewPostgreSQL(cfg *config.DatabaseConfig, log *logger.Logger) (*PostgreSQL, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)
	db.SetConnMaxIdleTime(30 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info().Msg("Connected to PostgreSQL database")

	return &PostgreSQL{DB: db}, nil
}

// Close closes the database connection
func (p *PostgreSQL) Close() error {
	if p.DB != nil {
		return p.DB.Close()
	}
	return nil
}

// Health checks if database is healthy
func (p *PostgreSQL) Health() error {
	return p.DB.Ping()
}
