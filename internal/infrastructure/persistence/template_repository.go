package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/hrygo/council/internal/core/workflow"
	"github.com/jackc/pgx/v5"
)

type TemplateRepository struct {
	pool DB
}

func NewTemplateRepository(pool DB) *TemplateRepository {
	return &TemplateRepository{pool: pool}
}

// Ensure table exists (Simplified migration for prototype)
func (r *TemplateRepository) Init(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS templates (
		id UUID PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		category TEXT NOT NULL DEFAULT 'custom',
		is_system BOOLEAN DEFAULT FALSE,
		graph JSONB NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);
	`
	_, err := r.pool.Exec(ctx, query)
	return err
}

func (r *TemplateRepository) List(ctx context.Context) ([]workflow.Template, error) {
	// Auto-init for now
	_ = r.Init(ctx)

	rows, err := r.pool.Query(ctx, "SELECT id, name, description, category, is_system, graph, created_at, updated_at FROM templates ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []workflow.Template
	for rows.Next() {
		var t workflow.Template
		var graphBytes []byte
		var cat string
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &cat, &t.IsSystem, &graphBytes, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		t.Category = workflow.TemplateCategory(cat)
		if err := json.Unmarshal(graphBytes, &t.Graph); err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}
	return templates, nil
}

func (r *TemplateRepository) Create(ctx context.Context, t *workflow.Template) error {
	_ = r.Init(ctx)

	graphJSON, err := json.Marshal(t.Graph)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO templates (id, name, description, category, is_system, graph, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = r.pool.Exec(ctx, query, t.ID, t.Name, t.Description, string(t.Category), t.IsSystem, graphJSON, time.Now(), time.Now())
	return err
}

func (r *TemplateRepository) Get(ctx context.Context, id string) (*workflow.Template, error) {
	_ = r.Init(ctx)

	var t workflow.Template
	var graphBytes []byte
	var cat string
	err := r.pool.QueryRow(ctx, "SELECT id, name, description, category, is_system, graph, created_at, updated_at FROM templates WHERE id = $1", id).
		Scan(&t.ID, &t.Name, &t.Description, &cat, &t.IsSystem, &graphBytes, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("template not found")
		}
		return nil, err
	}
	t.Category = workflow.TemplateCategory(cat)
	if err := json.Unmarshal(graphBytes, &t.Graph); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TemplateRepository) Delete(ctx context.Context, id string) error {
	_ = r.Init(ctx)
	_, err := r.pool.Exec(ctx, "DELETE FROM templates WHERE id = $1", id)
	return err
}
