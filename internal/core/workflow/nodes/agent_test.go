package nodes

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/hrygo/council/internal/core/agent"
	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/infrastructure/llm"
	"github.com/hrygo/council/internal/infrastructure/mocks"
)

func TestAgentProcessor_Process(t *testing.T) {
	// Setup Mocks
	mockRepo := mocks.NewAgentMockRepository()
	mockLLM := llm.NewMockProvider()
	mockLLM.StreamContent = []string{"Agent", " ", "Says", " ", "Hi"}

	agentID := uuid.New()
	mockRepo.Create(context.Background(), &agent.Agent{
		ID:            agentID,
		Name:          "TestAgent",
		PersonaPrompt: "You are a test agent.",
		ModelConfig:   agent.ModelConfig{Model: "gpt-4"},
	})

	processor := &AgentProcessor{
		AgentID:   agentID.String(),
		AgentRepo: mockRepo,
		LLM:       mockLLM,
	}

	// Execute
	stream := make(chan workflow.StreamEvent, 100)
	input := map[string]interface{}{"task": "Do something"}

	output, err := processor.Process(context.Background(), input, stream)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	// Verify
	if out, ok := output["agent_output"].(string); !ok || out != "Agent Says Hi" {
		t.Errorf("Expected 'Agent Says Hi', got '%v'", out)
	}
}
