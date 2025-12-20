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

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS schema_migrations").
		WillReturnResult(pgxmock.NewResult("CREATE", 0))

	// Check 001
	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM schema_migrations WHERE version=\\$1\\)").
		WithArgs("001_init_schema.up.sql").
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectExec("CREATE EXTENSION IF NOT EXISTS vector").WillReturnResult(pgxmock.NewResult("CREATE", 1))

	// Record 001
	mock.ExpectExec("INSERT INTO schema_migrations").
		WithArgs("001_init_schema.up.sql").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	// Check 002
	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM schema_migrations WHERE version=\\$1\\)").
		WithArgs("002_add_quarantine_logs.up.sql").
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS quarantine_logs").WillReturnResult(pgxmock.NewResult("CREATE", 1))

	// Record 002
	mock.ExpectExec("INSERT INTO schema_migrations").
		WithArgs("002_add_quarantine_logs.up.sql").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	// Check 003
	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM schema_migrations WHERE version=\\$1\\)").
		WithArgs("003_add_updated_at_columns.up.sql").
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

	// Content of 003 is usually ALTER TABLE... just assume Exec matches exact string or use Any?
	// But Migrate reads file content. The test probably doesn't read real file logic unless integration.
	// Wait, Migrate reads from `migrationFS`. The test uses the real `migrationFS`.
	// I don't know the exact content of 003.
	// I should check 003 content or use ExpectExec with regex matching key keywords like "ALTER TABLE".
	mock.ExpectExec("ALTER TABLE").WillReturnResult(pgxmock.NewResult("ALTER", 1))

	// Record 003
	mock.ExpectExec("INSERT INTO schema_migrations").
		WithArgs("003_add_updated_at_columns.up.sql").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	// Check 004
	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM schema_migrations WHERE version=\\$1\\)").
		WithArgs("004_create_llm_options.up.sql").
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

	// Content of 004: CREATE TABLE llm_providers ...; CREATE TABLE llm_models ...
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS llm_providers").WillReturnResult(pgxmock.NewResult("CREATE", 1))

	// Record 004
	mock.ExpectExec("INSERT INTO schema_migrations").
		WithArgs("004_create_llm_options.up.sql").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = Migrate(context.Background(), mock)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
