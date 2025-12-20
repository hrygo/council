package middleware

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/hrygo/council/internal/core/workflow"
)

// VersioningMiddleware creates automatic backups before HumanReview execution.
// This implements SPEC-605 to support the rollback functionality.
type VersioningMiddleware struct {
	BackupDir string
}

// NewVersioningMiddleware creates a new VersioningMiddleware instance.
func NewVersioningMiddleware(backupDir string) *VersioningMiddleware {
	return &VersioningMiddleware{BackupDir: backupDir}
}

// Name returns the middleware name.
func (v *VersioningMiddleware) Name() string {
	return "versioning"
}

// BeforeNodeExecution creates a backup before HumanReview node execution.
func (v *VersioningMiddleware) BeforeNodeExecution(
	ctx context.Context,
	session *workflow.Session,
	node *workflow.Node,
) error {
	// Only trigger for HumanReview nodes
	if node.Type != workflow.NodeTypeHumanReview {
		return nil
	}

	// Get target file from session inputs
	targetPath, ok := session.Inputs["target_file"].(string)
	if !ok || targetPath == "" {
		// No target file specified, skip backup
		return nil
	}

	// Create backup (log errors but don't block workflow)
	if err := v.CreateBackup(session.ID, targetPath); err != nil {
		// Log warning but don't fail the workflow
		fmt.Printf("[versioning] backup warning: %v\n", err)
	}

	return nil
}

// AfterNodeExecution is a no-op for versioning middleware.
func (v *VersioningMiddleware) AfterNodeExecution(
	ctx context.Context,
	session *workflow.Session,
	node *workflow.Node,
	output map[string]interface{},
) (map[string]interface{}, error) {
	return output, nil
}

// CreateBackup creates a timestamped backup of the target file.
func (v *VersioningMiddleware) CreateBackup(sessionID, targetPath string) error {
	// Check if source file exists
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return fmt.Errorf("target file does not exist: %s", targetPath)
	}

	// Create backup directory
	backupDir := filepath.Join(v.BackupDir, sessionID)
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Generate backup filename with timestamp
	filename := filepath.Base(targetPath)
	timestamp := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(backupDir, fmt.Sprintf("%s_%s.bak", filename, timestamp))

	// Copy file
	src, err := os.Open(targetPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	return nil
}

// FindLatestBackup returns the path to the most recent backup file.
func (v *VersioningMiddleware) FindLatestBackup(sessionID, targetPath string) string {
	filename := filepath.Base(targetPath)
	pattern := filepath.Join(v.BackupDir, sessionID, filename+"_*.bak")

	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		return ""
	}

	// Sort by filename (which includes timestamp) descending
	sort.Sort(sort.Reverse(sort.StringSlice(matches)))
	return matches[0]
}

// RestoreFromBackup copies backup content to original file.
func (v *VersioningMiddleware) RestoreFromBackup(backupPath, targetPath string) error {
	src, err := os.Open(backupPath)
	if err != nil {
		return fmt.Errorf("failed to open backup: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("failed to create target: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to restore file: %w", err)
	}

	return nil
}

// ListBackups returns all backup files for a session and target file.
func (v *VersioningMiddleware) ListBackups(sessionID, targetPath string) ([]string, error) {
	filename := filepath.Base(targetPath)
	pattern := filepath.Join(v.BackupDir, sessionID, filename+"_*.bak")

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	// Sort by timestamp descending (newest first)
	sort.Sort(sort.Reverse(sort.StringSlice(matches)))
	return matches, nil
}
