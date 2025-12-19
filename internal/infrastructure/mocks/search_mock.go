package mocks

import (
	"context"

	"github.com/hrygo/council/internal/infrastructure/search"
)

// SearchMockClient implements SearchClient for testing.
type SearchMockClient struct {
	Result *search.SearchResult
	Err    error
}

func (m *SearchMockClient) Search(ctx context.Context, query string, opts search.SearchOptions) (*search.SearchResult, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Result, nil
}
