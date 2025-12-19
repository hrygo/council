package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hrygo/council/internal/core/agent"
)

// AgentRepository implements agent.Repository using Postgres.
type AgentRepository struct {
	pool DB
}

// NewAgentRepository creates a new AgentRepository.
func NewAgentRepository(pool DB) *AgentRepository {
	return &AgentRepository{pool: pool}
}

func (r *AgentRepository) Create(ctx context.Context, a *agent.Agent) error {
	query := `
		INSERT INTO agents (name, avatar, description, persona_prompt, model_config, capabilities, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $7)
		RETURNING id, created_at, updated_at
	`
	err := r.pool.QueryRow(ctx, query,
		a.Name,
		a.Avatar,
		a.Description,
		a.PersonaPrompt,
		a.ModelConfig,
		a.Capabilities,
		time.Now(),
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)

	return err
}

func (r *AgentRepository) GetByID(ctx context.Context, id uuid.UUID) (*agent.Agent, error) {
	query := `SELECT id, name, avatar, description, persona_prompt, model_config, capabilities, created_at, updated_at FROM agents WHERE id = $1`
	var a agent.Agent
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&a.ID,
		&a.Name,
		&a.Avatar,
		&a.Description,
		&a.PersonaPrompt,
		&a.ModelConfig,
		&a.Capabilities,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent: %w", err)
	}
	return &a, nil
}

func (r *AgentRepository) List(ctx context.Context) ([]*agent.Agent, error) {
	query := `SELECT id, name, avatar, description, persona_prompt, model_config, capabilities, created_at, updated_at FROM agents ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []*agent.Agent
	for rows.Next() {
		var a agent.Agent
		if err := rows.Scan(
			&a.ID,
			&a.Name,
			&a.Avatar,
			&a.Description,
			&a.PersonaPrompt,
			&a.ModelConfig,
			&a.Capabilities,
			&a.CreatedAt,
			&a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		agents = append(agents, &a)
	}
	return agents, nil
}

func (r *AgentRepository) Update(ctx context.Context, a *agent.Agent) error {
	query := `
		UPDATE agents
		SET name = $1, avatar = $2, description = $3, persona_prompt = $4, model_config = $5, capabilities = $6, updated_at = $7
		WHERE id = $8
	`
	a.UpdatedAt = time.Now()
	_, err := r.pool.Exec(ctx, query,
		a.Name,
		a.Avatar,
		a.Description,
		a.PersonaPrompt,
		a.ModelConfig,
		a.Capabilities,
		a.UpdatedAt,
		a.ID,
	)
	return err
}

func (r *AgentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM agents WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
