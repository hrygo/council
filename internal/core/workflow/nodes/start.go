package nodes

import (
	"context"
	"time"

	"github.com/hrygo/council/internal/core/workflow"
)

type StartProcessor struct {
	OutputKeys []string
}

func (s *StartProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- workflow.StreamEvent) (map[string]interface{}, error) {
	// 1. Build Output from Input + Configured Keys
	output := map[string]interface{}{
		"metadata": map[string]interface{}{
			"started_at": time.Now(),
		},
	}

	workflow.ApplyPassthrough(input, output, workflow.PassthroughConfig{
		Keys: s.OutputKeys,
	})

	return output, nil
}
