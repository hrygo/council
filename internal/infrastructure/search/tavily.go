package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// SearchClient is the interface for web search providers.
type SearchClient interface {
	Search(ctx context.Context, query string, opts SearchOptions) (*SearchResult, error)
}

// SearchOptions configures the search request.
type SearchOptions struct {
	MaxResults int
	SearchType string // "search" | "answer"
	Domains    []string
}

// SearchResult holds the response from a search query.
type SearchResult struct {
	Query   string
	Results []SearchItem
	Answer  string // Direct answer (Tavily)
}

// SearchItem represents a single search result.
type SearchItem struct {
	Title   string
	URL     string
	Content string
	Score   float64
}

// TavilyClient implements SearchClient using the Tavily API.
type TavilyClient struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

// tavilyRequest is the request body for Tavily API.
type tavilyRequest struct {
	APIKey        string   `json:"api_key"`
	Query         string   `json:"query"`
	SearchDepth   string   `json:"search_depth"`
	IncludeAnswer bool     `json:"include_answer"`
	MaxResults    int      `json:"max_results"`
	Domains       []string `json:"include_domains,omitempty"`
}

// tavilyResponse is the response from Tavily API.
type tavilyResponse struct {
	Answer  string `json:"answer"`
	Results []struct {
		Title   string  `json:"title"`
		URL     string  `json:"url"`
		Content string  `json:"content"`
		Score   float64 `json:"score"`
	} `json:"results"`
}

// NewTavilyClient creates a new Tavily search client.
func NewTavilyClient() *TavilyClient {
	return &TavilyClient{
		APIKey:  os.Getenv("TAVILY_API_KEY"),
		BaseURL: "https://api.tavily.com",
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Search performs a web search using Tavily.
func (t *TavilyClient) Search(ctx context.Context, query string, opts SearchOptions) (*SearchResult, error) {
	if t.APIKey == "" {
		return nil, fmt.Errorf("TAVILY_API_KEY not set")
	}

	maxResults := opts.MaxResults
	if maxResults == 0 {
		maxResults = 5
	}

	reqBody := tavilyRequest{
		APIKey:        t.APIKey,
		Query:         query,
		SearchDepth:   "advanced",
		IncludeAnswer: true,
		MaxResults:    maxResults,
		Domains:       opts.Domains,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", t.BaseURL+"/search", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tavily request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("tavily returned status %d: %s", resp.StatusCode, string(body))
	}

	var tavilyResp tavilyResponse
	if err := json.NewDecoder(resp.Body).Decode(&tavilyResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	result := &SearchResult{
		Query:   query,
		Answer:  tavilyResp.Answer,
		Results: make([]SearchItem, len(tavilyResp.Results)),
	}

	for i, r := range tavilyResp.Results {
		result.Results[i] = SearchItem{
			Title:   r.Title,
			URL:     r.URL,
			Content: r.Content,
			Score:   r.Score,
		}
	}

	return result, nil
}

// NewSearchClient is a factory function for creating search clients.
func NewSearchClient(provider string) (SearchClient, error) {
	switch provider {
	case "tavily", "":
		return NewTavilyClient(), nil
	default:
		return nil, fmt.Errorf("unknown search provider: %s", provider)
	}
}
