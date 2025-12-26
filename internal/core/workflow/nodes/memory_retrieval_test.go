package nodes

import (
	"context"
	"testing"

	"github.com/hrygo/council/internal/core/memory"
	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/infrastructure/mocks"
	"github.com/stretchr/testify/assert"
)

func TestMemoryRetrievalProcessor_Process(t *testing.T) {
	mockManager := &mocks.MemoryMockManager{}
	processor := NewMemoryRetrievalProcessor(mockManager)

	t.Run("Successful retrieval", func(t *testing.T) {
		mockManager.RetrieveResult = []memory.ContextItem{
			{Content: "Historical context 1", Score: 0.9},
			{Content: "Historical context 2", Score: 0.8},
		}

		input := map[string]interface{}{
			"topic":    "AI Ethics",
			"group_id": "test-group",
		}
		stream := make(chan workflow.StreamEvent, 10)

		output, err := processor.Process(context.Background(), input, stream)
		assert.NoError(t, err)
		assert.Contains(t, output["history_context"].(string), "历史上下文")
		assert.Contains(t, output["history_context"].(string), "Historical context 1")
		assert.Contains(t, output["history_context"].(string), "Historical context 2")

		// Check stream events
		events := []workflow.StreamEvent{}
		for len(stream) > 0 {
			events = append(events, <-stream)
		}
		assert.Len(t, events, 2)
		assert.Equal(t, "node_state_change", events[0].Type)
		assert.Equal(t, "running", events[0].Data["status"])
		assert.Equal(t, "node_state_change", events[1].Type)
		assert.Equal(t, "completed", events[1].Data["status"])
	})

	t.Run("Empty topic", func(t *testing.T) {
		input := map[string]interface{}{
			"group_id": "test-group",
		}
		stream := make(chan workflow.StreamEvent, 10)

		output, err := processor.Process(context.Background(), input, stream)
		assert.NoError(t, err)
		assert.Empty(t, output["history_context"])
	})

	t.Run("Retrieval error", func(t *testing.T) {
		mockManager.Err = assert.AnError
		mockManager.RetrieveResult = nil

		input := map[string]interface{}{
			"topic":    "AI Ethics",
			"group_id": "test-group",
		}
		stream := make(chan workflow.StreamEvent, 10)

		output, err := processor.Process(context.Background(), input, stream)
		assert.NoError(t, err) // Should not fail according to code
		assert.Empty(t, output["history_context"])

		// Check error event
		foundError := false
		for len(stream) > 0 {
			event := <-stream
			if event.Type == "memory_retrieval_error" {
				foundError = true
				break
			}
		}
		assert.True(t, foundError)
	})
}
