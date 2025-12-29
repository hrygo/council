package nodes

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hrygo/council/internal/core/workflow"
)

type ContextSynthesizerProcessor struct {
	MaxRecentRounds int
}

func (p *ContextSynthesizerProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- workflow.StreamEvent) (map[string]interface{}, error) {
	stream <- workflow.StreamEvent{
		Type:      "node_state_change",
		Timestamp: time.Now(),
		Data:      map[string]interface{}{"status": "running"},
	}

	// Inputs
	currentSummary, _ := input["history_summary"].(string)
	newVerdict, _ := input["new_verdict"].(string)
	roundSummary, _ := input["round_summary"].(string) // Optional one-liner input

	// 1. Parse Existing
	legacy, rounds := parseContext(currentSummary)

	// 2. Append New
	// Inject hidden summary for future pruning
	fullTextWithMeta := newVerdict
	if roundSummary != "" {
		fullTextWithMeta += fmt.Sprintf("\n\n<!-- S: %s -->", roundSummary)
	}

	rounds = append(rounds, roundData{
		FullText: fullTextWithMeta,
		Summary:  roundSummary,
	})

	// 3. Prune
	if len(rounds) > p.MaxRecentRounds {
		toPrune := len(rounds) - p.MaxRecentRounds
		for i := 0; i < toPrune; i++ {
			r := rounds[i]
			// Add to legacy
			if r.Summary != "" {
				legacy = append(legacy, r.Summary)
			} else {
				// Fallback if no summary provided (simple truncation or extract first line)
				lines := strings.Split(r.FullText, "\n")
				if len(lines) > 0 {
					legacy = append(legacy, lines[0]) // Very naive
				}
			}
		}
		// Keep remaining
		rounds = rounds[toPrune:]
	}

	// 4. Reconstruct
	var sb strings.Builder
	if len(legacy) > 0 {
		sb.WriteString("## Legacy Context\n")
		for _, l := range legacy {
			sb.WriteString(fmt.Sprintf("- %s\n", l))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## Chronological Verdicts\n")
	for _, r := range rounds {
		sb.WriteString(r.FullText)
		sb.WriteString("\n\n")
	}

	final := sb.String()

	output := map[string]interface{}{
		"history_summary": final,
		"timestamp":       time.Now(),
	}

	stream <- workflow.StreamEvent{
		Type:      "node_state_change",
		Timestamp: time.Now(),
		Data:      map[string]interface{}{"status": "completed"},
	}

	return output, nil
}

type roundData struct {
	FullText string
	Summary  string
}

// Naive parser for TDD
func parseContext(text string) ([]string, []roundData) {
	var legacy []string
	var rounds []roundData

	if text == "" {
		return legacy, rounds
	}

	parts := strings.Split(text, "## Chronological Verdicts")

	// Parse Legacy
	if len(parts) > 0 {
		legacyPart := parts[0]
		if strings.Contains(legacyPart, "## Legacy Context") {
			lines := strings.Split(legacyPart, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "- ") {
					legacy = append(legacy, strings.TrimPrefix(line, "- "))
				}
			}
		}
	}

	// Parse Rounds (Assume separated by \n\n or some heuristic, but for now strict)
	// SPEC didn't define delimiter strictness. Let's assume input text blocks are distinct.
	// But `currentSummary` is a plain string. We need to split it back?
	// This is hard to do robustly without structured storage.
	// But SPEC-1104 says "Inputs: HistorySummary". It implies we store state as text?
	// Yes, for LLM Context input.

	// RE-DESIGN Implementation:
	// Storing strict structure is better. But if we must operate on Text Blob:
	// We need a specific delimiter for rounds.

	// Let's assume rounds involve "## Round" headers.

	// Wait, let's fix parsing based on provided `newVerdict` inputs in test.
	// Test inputs start with "## Round X".
	// We can split by "## Round".

	if len(parts) > 1 {
		verdictsPart := parts[1]
		rawRounds := strings.Split(verdictsPart, "## Round")
		for _, rr := range rawRounds {
			if strings.TrimSpace(rr) == "" {
				continue
			}
			full := "## Round" + rr
			// extract hidden summary?
			summary := extractSummary(full)
			rounds = append(rounds, roundData{
				FullText: strings.TrimSpace(full),
				Summary:  summary,
			})
		}
	}

	return legacy, rounds
}

func extractSummary(text string) string {
	// Look for special marker
	start := strings.Index(text, "<!-- S:")
	if start != -1 {
		end := strings.Index(text[start:], "-->")
		if end != -1 {
			return strings.TrimSpace(text[start+7 : start+end])
		}
	}
	// Fallback: Use title
	lines := strings.Split(text, "\n")
	if len(lines) > 0 {
		return strings.TrimLeft(lines[0], "# ")
	}
	return ""
}
