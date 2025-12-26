package db

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// Migrate runs embedded SQL migrations.
// Note: This is a simplified migrator for MVP. For production, consider using golang-migrate.
func Migrate(ctx context.Context, pool DB) error {
	// Read migration files
	files, err := fs.ReadDir(migrationFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations dir: %w", err)
	}

	// Filter for .up.sql files
	var ups []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".up.sql") {
			ups = append(ups, f.Name())
		}
	}
	sort.Strings(ups)

	log.Printf("Found %d migrations to apply", len(ups))

	// Ensure schema_migrations table exists
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ DEFAULT NOW()
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to ensure schema_migrations table: %w", err)
	}

	// Execute each migration
	for _, name := range ups {
		// Check if already applied
		var exists bool
		err := pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version=$1)", name).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check migration status for %s: %w", name, err)
		}

		if exists {
			// Skip quietly
			continue
		}

		content, err := migrationFS.ReadFile(filepath.Join("migrations", name))
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", name, err)
		}

		log.Printf("Applying migration: %s", name)
		if _, err := pool.Exec(ctx, string(content)); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", name, err)
		}

		// Record as applied
		_, err = pool.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", name)
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %w", name, err)
		}
	}

	return nil
}
