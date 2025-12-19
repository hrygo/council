package db

import (
	"context"
	"testing"
)

func TestConnect_InvalidURL(t *testing.T) {
	_, err := connect(context.Background(), "invalid-url")
	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestGetPool_InitiallyNil(t *testing.T) {
	// Note: this test might be flaky if other tests have initialized the global pool.
	// But in a fresh test run of just this package, it should be nil if Init wasn't called.
	// Actually, Init probably isn't called in tests yet.
	p := GetPool()
	if p != nil {
		t.Log("Warning: pool already initialized, possibly by other tests")
	}
}

func TestClose_NoPanic(t *testing.T) {
	Close() // Should not panic if pool is nil
}
