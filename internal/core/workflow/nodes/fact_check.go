package nodes

import (
	"context"
	"fmt"
	"time"

	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/infrastructure/llm"
)

type FactCheckProcessor struct {
	LLM             llm.LLMProvider
	SearchSources   []string
	VerifyThreshold float64
}

func (f *FactCheckProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- workflow.StreamEvent) (map[string]interface{}, error) {
	stream <- workflow.StreamEvent{
		Type:      "node_state_change",
		Timestamp: time.Now(),
		Data:      map[string]interface{}{"status": "running"},
	}

	// 1. Gather text to check from input
	var textToCheck string
	// Simplified extraction
	for _, v := range input {
		if s, ok := v.(string); ok {
			textToCheck += s + "\n"
		}
	}

	// 2. Perform Search (Mocked for now)
	// In real implementation, use f.SearchSources to call SerpAPI/Tavily
	searchResults := "Search results: [Fact 1 confirmed], [Fact 2 disputed]" // Placeholder

	// 3. Verify with LLM
	prompt := fmt.Sprintf(`Analyze the following text against search results.
Text: %s
Search Results: %s
Output 'VERIFIED' if accurate, 'DISPUTED' if false. Provide confidence 0.0-1.0.`, textToCheck, searchResults)

	// Call LLM... (Simplified synchronous call or reuse Stream logic from Agent)
	// For this sprint part, we'll simulate output.
	// In the future: f.LLM.Complete(ctx, prompt)
	_ = prompt // Silence unused variable error for now

	verified := true
	confidence := 0.9

	output := map[string]interface{}{
		"verified":   verified,
		"confidence": confidence,
		"correction": "None needed",
		"timestamp":  time.Now(),
	}

	stream <- workflow.StreamEvent{
		Type:      "node_state_change",
		Timestamp: time.Now(),
		Data:      map[string]interface{}{"status": "completed", "verified": verified},
	}

	return output, nil
}
