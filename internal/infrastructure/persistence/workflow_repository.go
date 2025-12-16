package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hrygo/council/internal/core/workflow"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkflowRepository struct {
	pool *pgxpool.Pool
}

func NewWorkflowRepository(pool *pgxpool.Pool) *WorkflowRepository {
	return &WorkflowRepository{pool: pool}
}

// WorkflowEntity represents the DB row
type WorkflowEntity struct {
	ID              string                   `json:"id"`
	GroupID         string                   `json:"group_id"`
	Name            string                   `json:"name"`
	GraphDefinition workflow.GraphDefinition `json:"graph_definition"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
}

func (r *WorkflowRepository) Create(ctx context.Context, graph *workflow.GraphDefinition) error {
	query := `
		INSERT INTO workflows (id, name, graph_definition)
		VALUES ($1, $2, $3)
	`
	// For MVP, we are not strictly enforcing GroupID yet in the input GraphDefinition
	// But the schema might require it if not nullable?
	// Schema: group_id UUID REFERENCES groups(id) ON DELETE CASCADE
	// It doesn't say NOT NULL. Let's check schema.
	// 001_init_schema.up.sql: group_id UUID REFERENCES groups(id) ON DELETE CASCADE
	// It allows NULL. So system workflows or global workflows can be null.

	if graph.ID == "" {
		graph.ID = uuid.New().String()
	}

	graphJSON, err := json.Marshal(graph)
	if err != nil {
		return fmt.Errorf("failed to marshal graph definition: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, graph.ID, graph.Name, graphJSON)
	if err != nil {
		return fmt.Errorf("failed to create workflow: %w", err)
	}
	return nil
}

func (r *WorkflowRepository) Get(ctx context.Context, id string) (*workflow.GraphDefinition, error) {
	query := `
		SELECT graph_definition FROM workflows WHERE id = $1
	`
	var graphJSON []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(&graphJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}

	var graph workflow.GraphDefinition
	if err := json.Unmarshal(graphJSON, &graph); err != nil {
		return nil, fmt.Errorf("failed to unmarshal graph definition: %w", err)
	}
	// Enforce ID match just in case JSON is stale, though it should be source of truth
	graph.ID = id
	return &graph, nil
}

func (r *WorkflowRepository) Update(ctx context.Context, graph *workflow.GraphDefinition) error {
	query := `
		UPDATE workflows 
		SET name = $2, graph_definition = $3, updated_at = NOW()
		WHERE id = $1
	`
	graphJSON, err := json.Marshal(graph)
	if err != nil {
		return fmt.Errorf("failed to marshal graph definition: %w", err)
	}

	cmdTag, err := r.pool.Exec(ctx, query, graph.ID, graph.Name, graphJSON)
	if err != nil {
		return fmt.Errorf("failed to update workflow: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("workflow not found")
	}
	return nil
}

func (r *WorkflowRepository) List(ctx context.Context) ([]*WorkflowEntity, error) {
	query := `
		SELECT id, name, graph_definition, created_at, updated_at FROM workflows ORDER BY updated_at DESC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list workflows: %w", err)
	}
	defer rows.Close()

	var list []*WorkflowEntity
	for rows.Next() {
		var w WorkflowEntity
		var graphJSON []byte
		if err := rows.Scan(&w.ID, &w.Name, &graphJSON, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(graphJSON, &w.GraphDefinition); err != nil {
			return nil, err
		}
		list = append(list, &w)
	}
	return list, nil
}
