package nodes

import (
	"context"
	"testing"

	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/infrastructure/llm"
	"github.com/hrygo/council/internal/infrastructure/mocks"
	"github.com/hrygo/council/internal/infrastructure/search"
)

func TestFactCheckProcessor_Process(t *testing.T) {
	// Setup Mocks
	mockLLM := llm.NewMockProvider()
	mockSearch := &mocks.SearchMockClient{
		Result: &search.SearchResult{
			Answer: "Mocked fact answer",
			Results: []search.SearchItem{
				{Title: "Source 1", Content: "Fact content"},
			},
		},
	}

	processor := &FactCheckProcessor{
		LLM:             mockLLM,
		SearchClient:    mockSearch,
		VerifyThreshold: 0.8,
	}

	// Case 1: All good (Verified = true)
	mockLLM.GenerateResponse = &llm.CompletionResponse{
		Content: `{"verified": true, "confidence": 0.95}`,
	}

	stream := make(chan workflow.StreamEvent, 10)
	input := map[string]interface{}{"text": "The sky is blue."}

	output, err := processor.Process(context.Background(), input, stream)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	if output["verified"] != true {
		t.Errorf("Expected verified=true, got %v", output["verified"])
	}

	// Case 2: Flagged (Verified = false)
	mockLLM.GenerateResponse = &llm.CompletionResponse{
		Content: `{"verified": false, "confidence": 0.5}`,
	}

	output, err = processor.Process(context.Background(), input, stream)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	if output["verified"] != false {
		t.Errorf("Expected verified=false, got %v", output["verified"])
	}
}
