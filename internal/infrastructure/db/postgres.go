package db

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
	once sync.Once
)

// Init initializes the database connection pool (singleton).
func Init(ctx context.Context, databaseURL string) error {
	var err error
	once.Do(func() {
		pool, err = connect(ctx, databaseURL)
	})
	return err
}

func connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	// Parse configuration
	config, parseErr := pgxpool.ParseConfig(databaseURL)
	if parseErr != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", parseErr)
	}

	// Create pool
	p, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Ping to verify connection
	if pingErr := p.Ping(ctx); pingErr != nil {
		p.Close()
		return nil, fmt.Errorf("failed to ping database: %w", pingErr)
	}

	// Run Migrations - using Wrapper to match Migrate signature change
	if migrateErr := Migrate(ctx, p); migrateErr != nil {
		p.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", migrateErr)
	}

	log.Println("Database connection established and migrations applied successfully")
	return p, nil
}

// GetPool returns the database connection pool.
func GetPool() *pgxpool.Pool {
	return pool
}

// Close closes the database connection pool.
func Close() {
	if pool != nil {
		pool.Close()
	}
}
