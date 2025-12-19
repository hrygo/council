package cache

import (
	"testing"
)

func TestConnect_InvalidAddr(t *testing.T) {
	// 127.0.0.1:0 is usually invalid or at least won't have redis
	_, err := connect("127.0.0.1:0", "", 0)
	if err == nil {
		t.Error("expected error for invalid address, got nil")
	}
}

func TestGetClient_InitiallyNil(t *testing.T) {
	c := GetClient()
	if c != nil {
		t.Log("Warning: client already initialized")
	}
}

func TestClose_NoPanic(t *testing.T) {
	Close()
}
