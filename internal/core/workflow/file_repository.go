package workflow

import (
	"context"
	"time"
)

// FileEntity represents a versioned file in the Virtual File System.
type FileEntity struct {
	ID        string    `json:"file_uuid" db:"file_uuid"`
	SessionID string    `json:"session_uuid" db:"session_uuid"`
	Path      string    `json:"path" db:"path"`
	Version   int       `json:"version" db:"version"`
	Content   string    `json:"content" db:"content"`
	Author    string    `json:"author" db:"author"`
	Reason    string    `json:"reason" db:"reason"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// SessionFileRepository defines the interface for VFS persistence.
type SessionFileRepository interface {
	// AddVersion creates a new version of a file.
	// It should handle concurrency by ensuring monotonic versioning or returning error on conflict.
	AddVersion(ctx context.Context, sessionID, path, content, author, reason string) (int, error)

	// GetLatest returns the highest version of a file.
	GetLatest(ctx context.Context, sessionID, path string) (*FileEntity, error)

	// ListFiles returns the latest version of all files in the session.
	ListFiles(ctx context.Context, sessionID string) (map[string]*FileEntity, error)

	// ListVersions returns all versions of a specific file.
	ListVersions(ctx context.Context, sessionID, path string) ([]*FileEntity, error)
}
