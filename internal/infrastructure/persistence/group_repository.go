package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hrygo/council/internal/core/group"
	"github.com/hrygo/council/internal/infrastructure/db"
)

// GroupRepository implements group.Repository using Postgres.
type GroupRepository struct {
	pool db.DB
}

// NewGroupRepository creates a new GroupRepository.
func NewGroupRepository(pool db.DB) *GroupRepository {
	return &GroupRepository{pool: pool}
}

func (r *GroupRepository) Create(ctx context.Context, g *group.Group) error {
	query := `
		INSERT INTO groups (name, icon, system_prompt, default_agent_uuids, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $5)
		RETURNING group_uuid, created_at, updated_at
	`
	// Ensure default_agent_uuids is empty array if nil, but DB default handles it.
	// However, we satisfy the query.
	if g.DefaultAgentUUIDs == nil {
		g.DefaultAgentUUIDs = []uuid.UUID{}
	}

	err := r.pool.QueryRow(ctx, query,
		g.Name,
		g.Icon,
		g.SystemPrompt,
		g.DefaultAgentUUIDs,
		time.Now(),
	).Scan(&g.ID, &g.CreatedAt, &g.UpdatedAt)

	return err
}

func (r *GroupRepository) GetByID(ctx context.Context, id uuid.UUID) (*group.Group, error) {
	query := `SELECT group_uuid, name, icon, system_prompt, default_agent_uuids, created_at, updated_at FROM groups WHERE group_uuid = $1`
	var g group.Group
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&g.ID,
		&g.Name,
		&g.Icon,
		&g.SystemPrompt,
		&g.DefaultAgentUUIDs,
		&g.CreatedAt,
		&g.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %w", err)
	}
	return &g, nil
}

func (r *GroupRepository) List(ctx context.Context) ([]*group.Group, error) {
	query := `SELECT group_uuid, name, icon, system_prompt, default_agent_uuids, created_at, updated_at FROM groups ORDER BY created_at DESC`
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
			&g.DefaultAgentUUIDs,
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
		SET name = $1, icon = $2, system_prompt = $3, default_agent_uuids = $4, updated_at = $5
		WHERE group_uuid = $6
	`
	g.UpdatedAt = time.Now()
	_, err := r.pool.Exec(ctx, query,
		g.Name,
		g.Icon,
		g.SystemPrompt,
		g.DefaultAgentUUIDs,
		g.UpdatedAt,
		g.ID,
	)
	return err
}

func (r *GroupRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM groups WHERE group_uuid = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
