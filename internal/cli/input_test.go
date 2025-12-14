package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewInputReader(t *testing.T) {
	var stdin, out bytes.Buffer
	reader := NewInputReader(&stdin, &out)

	if reader == nil {
		t.Fatal("NewInputReader() returned nil")
	}
	if reader.stdin != &stdin {
		t.Error("NewInputReader() did not set stdin correctly")
	}
	if reader.out != &out {
		t.Error("NewInputReader() did not set out correctly")
	}
}

func TestDefaultInputReader(t *testing.T) {
	reader := DefaultInputReader()
	if reader == nil {
		t.Fatal("DefaultInputReader() returned nil")
	}
}

func TestInputReader_ReadFile(t *testing.T) {
	// Create a temporary file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	content := "test file content"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	reader := DefaultInputReader()
	result, err := reader.ReadFile(tmpFile)

	if err != nil {
		t.Errorf("ReadFile() error = %v", err)
	}
	if result != content {
		t.Errorf("ReadFile() = %q, want %q", result, content)
	}
}

func TestInputReader_ReadFile_NotExists(t *testing.T) {
	reader := DefaultInputReader()
	_, err := reader.ReadFile("/nonexistent/file.txt")

	if err == nil {
		t.Error("ReadFile() should return error for non-existent file")
	}
}

func TestInputReader_ReadStdin(t *testing.T) {
	stdin := strings.NewReader("stdin content")
	reader := NewInputReader(stdin, &bytes.Buffer{})

	result, err := reader.ReadStdin()

	if err != nil {
		t.Errorf("ReadStdin() error = %v", err)
	}
	if result != "stdin content" {
		t.Errorf("ReadStdin() = %q, want %q", result, "stdin content")
	}
}

func TestInputReader_ReadInteractive(t *testing.T) {
	// Simulate user input with two empty lines at the end
	input := "line 1\nline 2\n\n\n"
	stdin := strings.NewReader(input)
	out := &bytes.Buffer{}
	reader := NewInputReader(stdin, out)

	result, err := reader.ReadInteractive()

	if err != nil {
		t.Errorf("ReadInteractive() error = %v", err)
	}

	// Should contain the content lines
	if !strings.Contains(result, "line 1") {
		t.Error("ReadInteractive() should contain 'line 1'")
	}
	if !strings.Contains(result, "line 2") {
		t.Error("ReadInteractive() should contain 'line 2'")
	}

	// Output should contain prompt
	if !strings.Contains(out.String(), "请输入待分析的材料") {
		t.Error("ReadInteractive() should print prompt")
	}
}

func TestInputReader_ReadInteractive_TrimTrailingEmptyLines(t *testing.T) {
	// Input with trailing empty lines before the double-empty terminator
	input := "content\n\n\n\n"
	stdin := strings.NewReader(input)
	reader := NewInputReader(stdin, &bytes.Buffer{})

	result, err := reader.ReadInteractive()

	if err != nil {
		t.Errorf("ReadInteractive() error = %v", err)
	}

	// Should not end with empty lines
	if strings.HasSuffix(result, "\n") && result != "content" {
		// Trim and check
		trimmed := strings.TrimRight(result, "\n")
		if trimmed != "content" {
			t.Errorf("ReadInteractive() = %q, should be 'content' after trimming", result)
		}
	}
}

func TestInputReader_ReadMaterial_Interactive(t *testing.T) {
	input := "interactive content\n\n\n"
	stdin := strings.NewReader(input)
	reader := NewInputReader(stdin, &bytes.Buffer{})

	result, err := reader.ReadMaterial("", true)

	if err != nil {
		t.Errorf("ReadMaterial() error = %v", err)
	}
	if !strings.Contains(result, "interactive content") {
		t.Error("ReadMaterial() with interactive should read from stdin")
	}
}

func TestInputReader_ReadMaterial_Stdin(t *testing.T) {
	stdin := strings.NewReader("piped content")
	reader := NewInputReader(stdin, &bytes.Buffer{})

	result, err := reader.ReadMaterial("-", false)

	if err != nil {
		t.Errorf("ReadMaterial() error = %v", err)
	}
	if result != "piped content" {
		t.Errorf("ReadMaterial() = %q, want %q", result, "piped content")
	}
}

func TestInputReader_ReadMaterial_File(t *testing.T) {
	// Create a temporary file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	content := "file content"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	reader := DefaultInputReader()
	result, err := reader.ReadMaterial(tmpFile, false)

	if err != nil {
		t.Errorf("ReadMaterial() error = %v", err)
	}
	if result != content {
		t.Errorf("ReadMaterial() = %q, want %q", result, content)
	}
}

func TestValidateMaterial(t *testing.T) {
	tests := []struct {
		name     string
		material string
		wantErr  bool
	}{
		{
			name:     "valid content",
			material: "some content",
			wantErr:  false,
		},
		{
			name:     "empty string",
			material: "",
			wantErr:  true,
		},
		{
			name:     "whitespace only",
			material: "   \t\n  ",
			wantErr:  true,
		},
		{
			name:     "content with whitespace",
			material: "  content  ",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMaterial(tt.material)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMaterial(%q) error = %v, wantErr %v", tt.material, err, tt.wantErr)
			}
		})
	}
}
