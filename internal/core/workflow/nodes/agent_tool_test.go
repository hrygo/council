package nodes_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hrygo/council/internal/core/agent"
	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/core/workflow/nodes"
	"github.com/hrygo/council/internal/core/workflow/nodes/tools"
	"github.com/hrygo/council/internal/infrastructure/llm"
	"github.com/hrygo/council/internal/infrastructure/mocks"
	"github.com/hrygo/council/internal/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestAgent_ToolCall_HappyPath(t *testing.T) {
	// Setup Mocks
	mockRepo := mocks.NewAgentMockRepository()
	mockLLM := llm.NewMockProvider()

	// Prepare Response Sequence
	// 1. LLM calls write_file
	resp1 := &llm.CompletionResponse{
		Content: "I will write the file.",
		ToolCalls: []llm.ToolCall{
			{
				ID:   "call_1",
				Type: "function",
				Function: llm.FunctionCall{
					Name:      "write_file",
					Arguments: `{"path": "main.go", "content": "package main"}`,
				},
			},
		},
		Usage: llm.Usage{TotalTokens: 20},
	}
	// 2. LLM confirms (after tool output injection)
	resp2 := &llm.CompletionResponse{
		Content: "Done.",
		Usage:   llm.Usage{TotalTokens: 10},
	}
	mockLLM.GenerateResponseQueue = []*llm.CompletionResponse{resp1, resp2}

	// Create Session & VFS
	session := workflow.NewSession(nil, nil)
	session.SetFileRepository(mocks.NewMockSessionFileRepository()) // We need a mock file repo

	// Agent Setup
	agentID := uuid.New()
	err := mockRepo.Create(context.Background(), &agent.Agent{
		ID:            agentID,
		Name:          "Surgeon",
		PersonaPrompt: "You are Surgeon.",
		ModelConfig:   agent.ModelConfig{Model: "gpt-4", Provider: "default"},
	})
	assert.NoError(t, err)

	registry := llm.NewRegistry(&config.Config{})
	registry.RegisterProvider("default", mockLLM)

	processor := &nodes.AgentProcessor{
		NodeID:    "surgeon",
		AgentID:   agentID.String(),
		AgentRepo: mockRepo,
		Registry:  registry,
		Tools:     []tools.Tool{&tools.WriteFileTool{}},
		Session:   session,
	}

	// Execute
	stream := make(chan workflow.StreamEvent, 100)
	input := map[string]interface{}{"task": "Create main.go"}

	output, err := processor.Process(context.Background(), input, stream)
	assert.NoError(t, err)

	// Verify Output
	assert.Equal(t, "Done.", output["agent_output"])

	// Verify VFS side effect
	f, err := session.GetLatestFile("main.go")
	assert.NoError(t, err)
	assert.Equal(t, "package main", f.Content)
	assert.Equal(t, 1, f.Version)
}
