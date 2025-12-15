package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Set test environment variables
	os.Setenv("PORT", "9090")
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	defer os.Unsetenv("PORT")
	defer os.Unsetenv("DATABASE_URL")

	cfg := Load()

	if cfg.Port != "9090" {
		t.Errorf("Expected port 9090, got %s", cfg.Port)
	}
	if cfg.DatabaseURL != "postgres://user:pass@localhost:5432/db" {
		t.Errorf("Expected db url, got %s", cfg.DatabaseURL)
	}

	// Verify defaults
	if cfg.Embedding.Provider != "siliconflow" {
		t.Errorf("Expected default embedding provider siliconflow, got %s", cfg.Embedding.Provider)
	}
}
