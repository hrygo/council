package workflow

import (
	"testing"
)

func TestGraphDefinition_Validate(t *testing.T) {
	tests := []struct {
		name    string
		graph   *GraphDefinition
		wantErr bool
	}{
		{
			name: "Valid Linear Graph",
			graph: &GraphDefinition{
				ID:          "linear",
				StartNodeID: "start",
				Nodes: map[string]*Node{
					"start": {ID: "start", Type: NodeTypeStart, NextIDs: []string{"step1"}},
					"step1": {ID: "step1", Type: NodeTypeAgent, NextIDs: []string{"end"}},
					"end":   {ID: "end", Type: NodeTypeEnd},
				},
			},
			wantErr: false,
		},
		{
			name: "Missing Start Node",
			graph: &GraphDefinition{
				ID:          "no_start",
				StartNodeID: "start",
				Nodes: map[string]*Node{
					"step1": {ID: "step1"},
				},
			},
			wantErr: true,
		},
		{
			name: "Node Pointing to Non-Existent Node",
			graph: &GraphDefinition{
				ID:          "broken_link",
				StartNodeID: "start",
				Nodes: map[string]*Node{
					"start": {ID: "start", Type: NodeTypeStart, NextIDs: []string{"ghost"}},
				},
			},
			wantErr: true,
		},
		{
			name: "Cycle Detected",
			graph: &GraphDefinition{
				ID:          "cycle",
				StartNodeID: "start",
				Nodes: map[string]*Node{
					"start": {ID: "start", Type: NodeTypeStart, NextIDs: []string{"step1"}},
					"step1": {ID: "step1", Type: NodeTypeAgent, NextIDs: []string{"step2"}},
					"step2": {ID: "step2", Type: NodeTypeAgent, NextIDs: []string{"step1"}}, // Cycle back to step1
				},
			},
			wantErr: true,
		},
		{
			name: "Isolated Node (Warning/Error logic depending on strictness - assuming strict for now)",
			graph: &GraphDefinition{
				ID:          "isolated",
				StartNodeID: "start",
				Nodes: map[string]*Node{
					"start":    {ID: "start", Type: NodeTypeStart, NextIDs: []string{"end"}},
					"end":      {ID: "end", Type: NodeTypeEnd},
					"isolated": {ID: "isolated", Type: NodeTypeAgent},
				},
			},
			wantErr: true, // Assuming we want to catch unreachable nodes
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.graph.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("GraphDefinition.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
