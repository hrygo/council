package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hrygo/council/internal/core/group"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GroupRepository implements group.Repository using Postgres.
type GroupRepository struct {
	pool *pgxpool.Pool
}

// NewGroupRepository creates a new GroupRepository.
func NewGroupRepository(pool *pgxpool.Pool) *GroupRepository {
	return &GroupRepository{pool: pool}
}

func (r *GroupRepository) Create(ctx context.Context, g *group.Group) error {
	query := `
		INSERT INTO groups (name, icon, system_prompt, default_agent_ids, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $5)
		RETURNING id, created_at, updated_at
	`
	// Ensure default_agent_ids is empty array if nil, but DB default handles it.
	// However, we satisfy the query.
	if g.DefaultAgentIDs == nil {
		g.DefaultAgentIDs = []uuid.UUID{}
	}

	err := r.pool.QueryRow(ctx, query,
		g.Name,
		g.Icon,
		g.SystemPrompt,
		g.DefaultAgentIDs,
		time.Now(),
	).Scan(&g.ID, &g.CreatedAt, &g.UpdatedAt)

	return err
}

func (r *GroupRepository) GetByID(ctx context.Context, id uuid.UUID) (*group.Group, error) {
	query := `SELECT id, name, icon, system_prompt, default_agent_ids, created_at, updated_at FROM groups WHERE id = $1`
	var g group.Group
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&g.ID,
		&g.Name,
		&g.Icon,
		&g.SystemPrompt,
		&g.DefaultAgentIDs,
		&g.CreatedAt,
		&g.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %w", err)
	}
	return &g, nil
}

func (r *GroupRepository) List(ctx context.Context) ([]*group.Group, error) {
	query := `SELECT id, name, icon, system_prompt, default_agent_ids, created_at, updated_at FROM groups ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*group.Group
	for rows.Next() {
		var g group.Group
		if err := rows.Scan(
			&g.ID,
			&g.Name,
			&g.Icon,
			&g.SystemPrompt,
			&g.DefaultAgentIDs,
			&g.CreatedAt,
			&g.UpdatedAt,
		); err != nil {
			return nil, err
		}
		groups = append(groups, &g)
	}
	return groups, nil
}

func (r *GroupRepository) Update(ctx context.Context, g *group.Group) error {
	query := `
		UPDATE groups
		SET name = $1, icon = $2, system_prompt = $3, default_agent_ids = $4, updated_at = $5
		WHERE id = $6
	`
	g.UpdatedAt = time.Now()
	_, err := r.pool.Exec(ctx, query,
		g.Name,
		g.Icon,
		g.SystemPrompt,
		g.DefaultAgentIDs,
		g.UpdatedAt,
		g.ID,
	)
	return err
}

func (r *GroupRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM groups WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
