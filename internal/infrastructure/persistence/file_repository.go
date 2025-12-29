package persistence

import (
	"context"
	"fmt"

	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/infrastructure/db"
)

type SessionFileRepository struct {
	pool db.DB
}

func NewSessionFileRepository(pool db.DB) workflow.SessionFileRepository {
	return &SessionFileRepository{pool: pool}
}

func (r *SessionFileRepository) AddVersion(ctx context.Context, sessionID, path, content, author, reason string) (int, error) {
	query := `
		INSERT INTO session_files (session_uuid, path, version, content, author, reason, created_at)
		VALUES ($1, $2, (
			SELECT COALESCE(MAX(version), 0) + 1 
			FROM session_files 
			WHERE session_uuid = $1 AND path = $2
		), $3, $4, $5, NOW())
		RETURNING version
	`
	var newVersion int
	err := r.pool.QueryRow(ctx, query, sessionID, path, content, author, reason).Scan(&newVersion)
	if err != nil {
		return 0, fmt.Errorf("failed to add file version: %w", err)
	}
	return newVersion, nil
}

func (r *SessionFileRepository) GetLatest(ctx context.Context, sessionID, path string) (*workflow.FileEntity, error) {
	query := `
		SELECT file_uuid, session_uuid, path, version, content, author, reason, created_at
		FROM session_files
		WHERE session_uuid = $1 AND path = $2
		ORDER BY version DESC
		LIMIT 1
	`
	var f workflow.FileEntity
	err := r.pool.QueryRow(ctx, query, sessionID, path).Scan(
		&f.ID, &f.SessionID, &f.Path, &f.Version, &f.Content, &f.Author, &f.Reason, &f.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest file version: %w", err)
	}
	return &f, nil
}

func (r *SessionFileRepository) ListFiles(ctx context.Context, sessionID string) (map[string]*workflow.FileEntity, error) {
	// We want the latest version for EACH path in the session.
	// DISTINCT ON (path) ORDER BY path, version DESC is the Postgres way.
	query := `
		SELECT DISTINCT ON (path) file_uuid, session_uuid, path, version, content, author, reason, created_at
		FROM session_files
		WHERE session_uuid = $1
		ORDER BY path, version DESC
	`
	rows, err := r.pool.Query(ctx, query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}
	defer rows.Close()

	files := make(map[string]*workflow.FileEntity)
	for rows.Next() {
		var f workflow.FileEntity
		if err := rows.Scan(&f.ID, &f.SessionID, &f.Path, &f.Version, &f.Content, &f.Author, &f.Reason, &f.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan file row: %w", err)
		}
		files[f.Path] = &f
	}
	return files, nil
}

func (r *SessionFileRepository) ListVersions(ctx context.Context, sessionID, path string) ([]*workflow.FileEntity, error) {
	query := `
		SELECT file_uuid, session_uuid, path, version, content, author, reason, created_at
		FROM session_files
		WHERE session_uuid = $1 AND path = $2
		ORDER BY version DESC
	`
	rows, err := r.pool.Query(ctx, query, sessionID, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list file versions: %w", err)
	}
	defer rows.Close()

	var versions []*workflow.FileEntity
	for rows.Next() {
		var f workflow.FileEntity
		if err := rows.Scan(&f.ID, &f.SessionID, &f.Path, &f.Version, &f.Content, &f.Author, &f.Reason, &f.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan file version row: %w", err)
		}
		versions = append(versions, &f)
	}
	return versions, nil
}
