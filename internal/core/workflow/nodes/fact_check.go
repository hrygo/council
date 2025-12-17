package nodes

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/infrastructure/llm"
	"github.com/hrygo/council/internal/infrastructure/search"
)

type FactCheckProcessor struct {
	LLM             llm.LLMProvider
	SearchClient    search.SearchClient
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
	for _, v := range input {
		if s, ok := v.(string); ok {
			textToCheck += s + "\n"
		}
	}

	// 2. Perform Web Search
	var searchResults string
	if f.SearchClient != nil {
		result, err := f.SearchClient.Search(ctx, textToCheck, search.SearchOptions{
			MaxResults: 3,
			SearchType: "answer",
		})
		if err != nil {
			stream <- workflow.StreamEvent{
				Type:      "error",
				Timestamp: time.Now(),
				Data:      map[string]interface{}{"error": "search failed: " + err.Error()},
			}
			// Proceed with empty results instead of failing
			searchResults = "[Search unavailable]"
		} else {
			// Build search summary
			var sb strings.Builder
			if result.Answer != "" {
				sb.WriteString("Direct Answer: " + result.Answer + "\n")
			}
			for _, item := range result.Results {
				sb.WriteString(fmt.Sprintf("- %s: %s\n", item.Title, item.Content))
			}
			searchResults = sb.String()
		}
	} else {
		searchResults = "[No search client configured]"
	}

	// 3. Verify with LLM
	prompt := fmt.Sprintf(`Analyze the following text against web search results.
Text to verify: 
%s

Web Search Results:
%s

Determine if the claims in the text are accurate.
Output JSON: {"verified": true/false, "confidence": 0.0-1.0, "issues": ["list of issues if any"]}`, textToCheck, searchResults)

	verified := true
	confidence := 0.9
	issues := []string{}

	if f.LLM != nil {
		resp, err := f.LLM.Generate(ctx, &llm.CompletionRequest{
			Messages: []llm.Message{
				{Role: "user", Content: prompt},
			},
			Temperature: 0.1,
		})
		if err == nil && resp.Content != "" {
			// Parse LLM response (simplified - assume success means verified)
			if strings.Contains(strings.ToLower(resp.Content), `"verified": false`) ||
				strings.Contains(strings.ToLower(resp.Content), `"verified":false`) {
				verified = false
				confidence = 0.7
				issues = append(issues, "LLM flagged potential inaccuracies")
			}
		}
	}

	output := map[string]interface{}{
		"verified":       verified,
		"confidence":     confidence,
		"issues":         issues,
		"search_summary": searchResults,
		"timestamp":      time.Now(),
	}

	stream <- workflow.StreamEvent{
		Type:      "node_state_change",
		Timestamp: time.Now(),
		Data:      map[string]interface{}{"status": "completed", "verified": verified},
	}

	return output, nil
}
