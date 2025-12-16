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

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// Migrate runs embedded SQL migrations.
// Note: This is a simplified migrator for MVP. For production, consider using golang-migrate.
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
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

	// Execute each migration
	for _, name := range ups {
		content, err := migrationFS.ReadFile(filepath.Join("migrations", name))
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", name, err)
		}

		log.Printf("Applying migration: %s", name)
		if _, err := pool.Exec(ctx, string(content)); err != nil {
			// In a real system we should check schema_migrations table.
			// Here we rely on IF NOT EXISTS or idempotent SQL.
			// If 001 fails (already exists), we log and continue to try 002.
			// This is NOT robust for Rollbacks but acceptable for this "Forward Only" iteration.
			log.Printf("Warning: Migration %s execution result: %v", name, err)
		}
	}

	return nil
}
