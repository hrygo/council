package nodes

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/hrygo/council/internal/core/agent"
	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/infrastructure/llm"
	"github.com/hrygo/council/internal/infrastructure/mocks"
	"github.com/hrygo/council/internal/pkg/config"
)

func TestAgentProcessor_Process(t *testing.T) {
	// Setup Mocks
	mockRepo := mocks.NewAgentMockRepository()
	mockLLM := llm.NewMockProvider()
	mockLLM.GenerateResponse = &llm.CompletionResponse{
		Content: "Agent Says Hi",
		Usage:   llm.Usage{TotalTokens: 10},
	}

	agentID := uuid.New()
	if err := mockRepo.Create(context.Background(), &agent.Agent{
		ID:            agentID,
		Name:          "TestAgent",
		PersonaPrompt: "You are a test agent.",
		ModelConfig:   agent.ModelConfig{Model: "gpt-4", Provider: "default"},
	}); err != nil {
		t.Fatalf("Failed to create mock agent: %v", err)
	}

	// Mock Registry
	cfg := &config.Config{}
	registry := llm.NewRegistry(cfg)
	registry.RegisterProvider("default", mockLLM)

	processor := &AgentProcessor{
		AgentID:   agentID.String(),
		AgentRepo: mockRepo,
		Registry:  registry,
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
