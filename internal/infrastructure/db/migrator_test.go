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

	// 001_init_schema.up.sql
	mock.ExpectExec("CREATE EXTENSION IF NOT EXISTS vector").WillReturnResult(pgxmock.NewResult("CREATE", 1))
	// 002_add_quarantine_logs.up.sql
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS quarantine_logs").WillReturnResult(pgxmock.NewResult("CREATE", 1))

	err = Migrate(context.Background(), mock)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
