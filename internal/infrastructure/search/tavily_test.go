package search

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTavilyClient_Search(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/search" {
			t.Errorf("expected /search, got %s", r.URL.Path)
		}

		resp := tavilyResponse{
			Answer: "Tavily answer",
			Results: []struct {
				Title   string  `json:"title"`
				URL     string  `json:"url"`
				Content string  `json:"content"`
				Score   float64 `json:"score"`
			}{
				{Title: "Result 1", URL: "http://example.com/1", Content: "Body 1", Score: 0.99},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	client := &TavilyClient{
		APIKey:  "test-key",
		BaseURL: ts.URL,
		Client:  ts.Client(),
	}

	res, err := client.Search(context.Background(), "test query", SearchOptions{MaxResults: 1})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Answer != "Tavily answer" {
		t.Errorf("expected 'Tavily answer', got '%s'", res.Answer)
	}
	if len(res.Results) != 1 || res.Results[0].Title != "Result 1" {
		t.Errorf("unexpected results: %+v", res.Results)
	}
}

func TestTavilyClient_NoAPIKey(t *testing.T) {
	client := &TavilyClient{APIKey: ""}
	_, err := client.Search(context.Background(), "q", SearchOptions{})
	if err == nil {
		t.Error("expected error for missing API key, got nil")
	}
}

func TestNewSearchClient(t *testing.T) {
	c, err := NewSearchClient("tavily")
	if err != nil || c == nil {
		t.Errorf("failed to create tavily client: %v", err)
	}

	_, err = NewSearchClient("unknown")
	if err == nil {
		t.Error("expected error for unknown provider, got nil")
	}
}
