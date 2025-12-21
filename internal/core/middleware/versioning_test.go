package middleware

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestVersioningMiddleware_CreateBackup(t *testing.T) {
	// Setup temp directories
	tmpDir := t.TempDir()
	sourceDir := filepath.Join(tmpDir, "source")
	backupDir := filepath.Join(tmpDir, "backup")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("failed to create source dir: %v", err)
	}

	// Create a test file
	testFile := filepath.Join(sourceDir, "test_doc.md")
	originalContent := "# Original Content\n\nThis is the original document."
	if err := os.WriteFile(testFile, []byte(originalContent), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Create middleware and backup
	mw := NewVersioningMiddleware(backupDir)
	sessionID := "test-session-123"

	if err := mw.CreateBackup(sessionID, testFile); err != nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}

	// Verify backup was created
	backups, err := mw.ListBackups(sessionID, testFile)
	if err != nil {
		t.Fatalf("ListBackups failed: %v", err)
	}
	if len(backups) != 1 {
		t.Errorf("expected 1 backup, got %d", len(backups))
	}

	// Verify backup content
	if len(backups) > 0 {
		content, _ := os.ReadFile(backups[0])
		if string(content) != originalContent {
			t.Errorf("backup content mismatch")
		}
	}
}

func TestVersioningMiddleware_FindLatestBackup(t *testing.T) {
	tmpDir := t.TempDir()
	backupDir := filepath.Join(tmpDir, "backup")
	sessionID := "session-456"
	sessionBackupDir := filepath.Join(backupDir, sessionID)
	if err := os.MkdirAll(sessionBackupDir, 0755); err != nil {
		t.Fatalf("failed to create session backup dir: %v", err)
	}

	// Create multiple backup files with different timestamps
	backupFiles := []string{
		"doc.md_20241220_100000.bak",
		"doc.md_20241220_120000.bak", // This should be latest
		"doc.md_20241220_110000.bak",
	}

	for _, f := range backupFiles {
		if err := os.WriteFile(filepath.Join(sessionBackupDir, f), []byte("content"), 0644); err != nil {
			t.Fatalf("failed to write test file %s: %v", f, err)
		}
	}

	mw := NewVersioningMiddleware(backupDir)
	latest := mw.FindLatestBackup(sessionID, "/some/path/doc.md")

	if !strings.HasSuffix(latest, "doc.md_20241220_120000.bak") {
		t.Errorf("expected latest backup ending with 20241220_120000.bak, got %s", latest)
	}
}

func TestVersioningMiddleware_RestoreFromBackup(t *testing.T) {
	tmpDir := t.TempDir()

	// Create backup file
	backupPath := filepath.Join(tmpDir, "backup.bak")
	backupContent := "This is the backed up content"
	if err := os.WriteFile(backupPath, []byte(backupContent), 0644); err != nil {
		t.Fatalf("failed to write backup file: %v", err)
	}

	// Create modified target file
	targetPath := filepath.Join(tmpDir, "target.md")
	if err := os.WriteFile(targetPath, []byte("Modified content"), 0644); err != nil {
		t.Fatalf("failed to write target file: %v", err)
	}

	// Restore
	mw := NewVersioningMiddleware(tmpDir)
	if err := mw.RestoreFromBackup(backupPath, targetPath); err != nil {
		t.Fatalf("RestoreFromBackup failed: %v", err)
	}

	// Verify restored content
	content, _ := os.ReadFile(targetPath)
	if string(content) != backupContent {
		t.Errorf("restore failed: expected '%s', got '%s'", backupContent, string(content))
	}
}

func TestVersioningMiddleware_CreateBackup_NonExistentFile(t *testing.T) {
	mw := NewVersioningMiddleware(t.TempDir())

	err := mw.CreateBackup("session", "/non/existent/file.md")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

func TestVersioningMiddleware_MultipleBackups(t *testing.T) {
	tmpDir := t.TempDir()
	sourceFile := filepath.Join(tmpDir, "doc.md")
	backupDir := filepath.Join(tmpDir, "backup")

	if err := os.WriteFile(sourceFile, []byte("v1"), 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	mw := NewVersioningMiddleware(backupDir)
	sessionID := "multi-backup-test"

	// Create multiple backups with different timestamps
	if err := mw.CreateBackup(sessionID, sourceFile); err != nil {
		t.Fatalf("first backup failed: %v", err)
	}
	time.Sleep(1100 * time.Millisecond) // Ensure different timestamps (at least 1 second apart)
	if err := os.WriteFile(sourceFile, []byte("v2"), 0644); err != nil {
		t.Fatalf("failed to update source file: %v", err)
	}
	if err := mw.CreateBackup(sessionID, sourceFile); err != nil {
		t.Fatalf("second backup failed: %v", err)
	}

	backups, _ := mw.ListBackups(sessionID, sourceFile)
	if len(backups) != 2 {
		t.Errorf("expected 2 backups, got %d", len(backups))
	}
}
