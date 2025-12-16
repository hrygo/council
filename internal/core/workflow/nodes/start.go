package nodes

import (
	"context"
	"strings"
	"time"

	"github.com/hrygo/council/internal/core/workflow"
)

type StartProcessor struct{}

func (s *StartProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- workflow.StreamEvent) (map[string]interface{}, error) {
	// 1. Notify Start
	stream <- workflow.StreamEvent{
		Type:      "node_state_change",
		Timestamp: time.Now(),
		Data:      map[string]interface{}{"node_id": "start", "status": "running"},
	}

	// 2. Parse Inputs
	proposal, _ := input["proposal"].(string)
	attachments, _ := input["attachments"].([]map[string]interface{})

	var parsedContents []string
	for _, att := range attachments {
		// Simplified for MVP: Assuming 'content' field is populated directly or via basic read
		if content, ok := att["content"].(string); ok {
			parsedContents = append(parsedContents, content)
		}
	}

	// 3. Construct Output
	output := map[string]interface{}{
		"proposal":         proposal,
		"attachments":      attachments,
		"combined_context": strings.Join(parsedContents, "\n\n---\n\n"),
		"metadata": map[string]interface{}{
			"started_at":       time.Now(),
			"attachment_count": len(attachments),
		},
	}

	// 4. Notify Complete
	stream <- workflow.StreamEvent{
		Type:      "node_state_change",
		Timestamp: time.Now(),
		Data:      map[string]interface{}{"node_id": "start", "status": "completed"},
	}

	return output, nil
}
