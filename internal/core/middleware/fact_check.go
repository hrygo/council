package middleware

import (
	"context"
	"regexp"

	"github.com/hrygo/council/internal/core/workflow"
)

// FactCheckTrigger scans output for unverified claims
type FactCheckTrigger struct {
	MetricRegex   *regexp.Regexp
	CitationRegex *regexp.Regexp
}

func NewFactCheckTrigger() *FactCheckTrigger {
	return &FactCheckTrigger{
		MetricRegex:   regexp.MustCompile(`\[Specific Metric\]`),
		CitationRegex: regexp.MustCompile(`\[External Citation\]`),
	}
}

func (fc *FactCheckTrigger) Name() string {
	return "AntiHallucination"
}

func (fc *FactCheckTrigger) BeforeNodeExecution(ctx context.Context, session *workflow.Session, node *workflow.Node) error {
	return nil
}

func (fc *FactCheckTrigger) AfterNodeExecution(ctx context.Context, session *workflow.Session, node *workflow.Node, output map[string]interface{}) (map[string]interface{}, error) {
	// Scan "content" field if exists
	if content, ok := output["content"].(string); ok {
		if fc.MetricRegex.MatchString(content) || fc.CitationRegex.MatchString(content) {
			// Found unverified claim.
			// Action: Inject system note or flag.
			// For MVP: We append a warning flag to metadata.
			// If output has "metadata", append.

			meta, _ := output["metadata"].(map[string]interface{})
			if meta == nil {
				meta = make(map[string]interface{})
			}
			meta["verify_pending"] = true
			output["metadata"] = meta

			// Optionally inject warning in content?
			// output["content"] = content + "\n\n[System Alert: Unverified claims detected. Fact Check required.]"
		}
	}
	return output, nil
}
