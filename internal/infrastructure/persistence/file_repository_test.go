package persistence

import (
	"context"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func TestAddVersion(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewSessionFileRepository(mock)
	ctx := context.Background()

	sessionID := "sess-123"
	path := "main.go"
	content := "package main"
	author := "user"
	reason := "init"

	mock.ExpectQuery("INSERT INTO session_files").
		WithArgs(sessionID, path, content, author, reason).
		WillReturnRows(pgxmock.NewRows([]string{"version"}).AddRow(1))

	ver, err := repo.AddVersion(ctx, sessionID, path, content, author, reason)
	assert.NoError(t, err)
	assert.Equal(t, 1, ver)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetLatest(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewSessionFileRepository(mock)
	ctx := context.Background()

	sessionID := "sess-123"
	path := "main.go"
	now := time.Now()

	mock.ExpectQuery("SELECT file_uuid, session_uuid, path, version, content, author, reason, created_at FROM session_files").
		WithArgs(sessionID, path).
		WillReturnRows(pgxmock.NewRows([]string{"file_uuid", "session_uuid", "path", "version", "content", "author", "reason", "created_at"}).
			AddRow("uuid-1", sessionID, path, 5, "content", "author", "reason", now))

	f, err := repo.GetLatest(ctx, sessionID, path)
	assert.NoError(t, err)
	assert.Equal(t, 5, f.Version)
	assert.Equal(t, "content", f.Content)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestListFiles(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewSessionFileRepository(mock)
	ctx := context.Background()

	sessionID := "sess-123"
	now := time.Now()

	mock.ExpectQuery("SELECT DISTINCT ON \\(path\\) .* FROM session_files").
		WithArgs(sessionID).
		WillReturnRows(pgxmock.NewRows([]string{"file_uuid", "session_uuid", "path", "version", "content", "author", "reason", "created_at"}).
			AddRow("uuid-1", sessionID, "main.go", 2, "pkg main", "agent", "fix", now).
			AddRow("uuid-2", sessionID, "README.md", 1, "# Hello", "user", "init", now))

	files, err := repo.ListFiles(ctx, sessionID)
	assert.NoError(t, err)
	assert.Len(t, files, 2)
	foundMain := false
	foundReadme := false
	for _, f := range files {
		if f.Path == "main.go" {
			assert.Equal(t, 2, f.Version)
			foundMain = true
		} else if f.Path == "README.md" {
			assert.Equal(t, 1, f.Version)
			foundReadme = true
		}
	}
	assert.True(t, foundMain, "main.go not found")
	assert.True(t, foundReadme, "README.md not found")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
