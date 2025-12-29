package llm

import (
	"encoding/json"
	"testing"
)

func TestMessageSerialization(t *testing.T) {
	msg := Message{
		Role:    "assistant",
		Content: "Call tool",
		ToolCalls: []ToolCall{
			{
				ID:   "call_123",
				Type: "function",
				Function: FunctionCall{
					Name:      "write_file",
					Arguments: `{"path": "main.go"}`,
				},
			},
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var parsed Message
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(parsed.ToolCalls) != 1 {
		t.Errorf("Expected 1 tool call, got %d", len(parsed.ToolCalls))
	}
	if parsed.ToolCalls[0].Function.Name != "write_file" {
		t.Errorf("Expected write_file, got %s", parsed.ToolCalls[0].Function.Name)
	}
}
