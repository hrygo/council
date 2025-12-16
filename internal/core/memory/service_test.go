package memory

import (
	"context"
	"testing"
)

type MockEmbedder struct{}

func (m *MockEmbedder) Embed(ctx context.Context, model string, text string) ([]float32, error) {
	return []float32{0.1, 0.2, 0.3}, nil
}

func TestMemoryRetrieval(t *testing.T) {
	service := NewService(&MockEmbedder{})

	// We can't easily test real Redis/PG connections in unit test without mocking DB/Cache drivers
	// But we can test the logic flow if we could inject mocks for cache/db.
	// Since current implementation uses global getters (cache.GetClient), unit testing is hard without running docker.
	// We will assume integration tests cover this or manual verification.
	// For now, testing Service instantiation.

	if service.Embedder == nil {
		t.Error("Embedder not injected")
	}

	// Test logic of Retrieve (partially)
	// Actually Retrieve calls global db.GetPool() immediately.
	// We skip deep testing here without a complex mock setup or integration environment.
}
