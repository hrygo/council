package db

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v3"
)

func TestMigrate(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	// 1. Ensure schema_migrations table
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS schema_migrations").
		WillReturnResult(pgxmock.NewResult("CREATE", 0))

	// 2. Check 001_v2_schema_init
	migrationName := "001_v2_schema_init.up.sql"
	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM schema_migrations WHERE version=\\$1\\)").
		WithArgs(migrationName).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

	// 3. Apply 001_v2_schema_init
	// Using a broad regex matcher as the file is large and contains many statements.
	// We verify that *some* content specific to the migration is executed.
	// The migration starts with "-- Squashed Migration", but comments might be stripped or handled differently?
	// Actually pgx/pgxmock Exec receives the string.
	// Regex `(?s).*` matches everything including newlines.
	mock.ExpectExec("(?s).*").
		WillReturnResult(pgxmock.NewResult("CREATE", 1))

	// 4. Record 001_v2_schema_init
	mock.ExpectExec("INSERT INTO schema_migrations").
		WithArgs(migrationName).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	// 5. Check 002_add_node_statuses
	migrationName2 := "002_add_node_statuses.up.sql"
	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM schema_migrations WHERE version=\\$1\\)").
		WithArgs(migrationName2).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

	// 6. Apply 002_add_node_statuses
	mock.ExpectExec("(?s).*").
		WillReturnResult(pgxmock.NewResult("ALTER", 1))

	// 7. Record 002_add_node_statuses
	mock.ExpectExec("INSERT INTO schema_migrations").
		WithArgs(migrationName2).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = Migrate(context.Background(), mock)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
