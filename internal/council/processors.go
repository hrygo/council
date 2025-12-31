package council

import (
	"context"
	"strings"
	"time"

	"github.com/hrygo/council/internal/core/workflow"
)

type StartProcessor struct {
	OutputKeys []string
}

func (s *StartProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- workflow.StreamEvent) (map[string]interface{}, error) {
	// 2. Base Output
	output := map[string]interface{}{
		"metadata": map[string]interface{}{
			"started_at": time.Now(),
		},
	}

	// 3. Council Logic: Parse attachments
	// Hande both inputs: []map[string]interface{} (from code) or []interface{} (from JSON)
	var attachments []interface{}
	if v, ok := input["attachments"].([]interface{}); ok {
		attachments = v
	} else if v, ok := input["attachments"].([]map[string]interface{}); ok {
		for _, item := range v {
			attachments = append(attachments, item)
		}
	}

	if len(attachments) > 0 {
		var parsedContents []string
		count := 0
		for _, rawAtt := range attachments {
			if att, ok := rawAtt.(map[string]interface{}); ok {
				if content, ok := att["content"].(string); ok {
					parsedContents = append(parsedContents, content)
					count++
				}
			}
		}
		if len(parsedContents) > 0 {
			output["combined_context"] = strings.Join(parsedContents, "\n\n---\n\n")
		}
		output["metadata"].(map[string]interface{})["attachment_count"] = count
	}

	// 4. Apply configured Passthrough
	workflow.ApplyPassthrough(input, output, workflow.PassthroughConfig{
		Keys: s.OutputKeys,
	})

	return output, nil
}
