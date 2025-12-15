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

// Init initializes the database connection pool.
func Init(ctx context.Context, databaseURL string) error {
	var err error
	once.Do(func() {
		// Parse configuration
		config, parseErr := pgxpool.ParseConfig(databaseURL)
		if parseErr != nil {
			err = fmt.Errorf("failed to parse database URL: %w", parseErr)
			return
		}

		// Create pool
		pool, err = pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			err = fmt.Errorf("failed to create connection pool: %w", err)
			return
		}

		// Ping to verify connection
		if pingErr := pool.Ping(ctx); pingErr != nil {
			err = fmt.Errorf("failed to ping database: %w", pingErr)
			return
		}

		log.Println("Database connection established successfully")
	})

	return err
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
