package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/infrastructure/db"
	"github.com/jackc/pgx/v5"
)

type TemplateRepository struct {
	pool db.DB
}

func NewTemplateRepository(pool db.DB) *TemplateRepository {
	return &TemplateRepository{pool: pool}
}

// Note: Init method removed as schema is managed by migrations.

func (r *TemplateRepository) List(ctx context.Context) ([]workflow.Template, error) {
	rows, err := r.pool.Query(ctx, "SELECT template_uuid, name, description, is_system, graph_definition, created_at, updated_at FROM workflow_templates ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []workflow.Template
	for rows.Next() {
		var t workflow.Template
		var graphBytes []byte
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.IsSystem, &graphBytes, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		// Default category based on IsSystem
		if t.IsSystem {
			t.Category = workflow.TemplateCategoryOther // Or specific logic
		} else {
			t.Category = workflow.TemplateCategoryCustom
		}

		if err := json.Unmarshal(graphBytes, &t.Graph); err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}
	return templates, nil
}

func (r *TemplateRepository) Create(ctx context.Context, t *workflow.Template) error {
	graphJSON, err := json.Marshal(t.Graph)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO workflow_templates (template_uuid, name, description, is_system, graph_definition, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	// Category is ignored as it's not in DB
	_, err = r.pool.Exec(ctx, query, t.ID, t.Name, t.Description, t.IsSystem, graphJSON, time.Now(), time.Now())
	return err
}

func (r *TemplateRepository) Get(ctx context.Context, id string) (*workflow.Template, error) {
	var t workflow.Template
	var graphBytes []byte
	err := r.pool.QueryRow(ctx, "SELECT template_uuid, name, description, is_system, graph_definition, created_at, updated_at FROM workflow_templates WHERE template_uuid = $1", id).
		Scan(&t.ID, &t.Name, &t.Description, &t.IsSystem, &graphBytes, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("template not found")
		}
		return nil, err
	}

	if t.IsSystem {
		t.Category = workflow.TemplateCategoryOther
	} else {
		t.Category = workflow.TemplateCategoryCustom
	}

	if err := json.Unmarshal(graphBytes, &t.Graph); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TemplateRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM workflow_templates WHERE template_uuid = $1", id)
	return err
}
