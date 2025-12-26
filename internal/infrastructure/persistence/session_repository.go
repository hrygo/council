package persistence

import (
	"context"
	"fmt"

	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/infrastructure/db"
)

type SessionRepository struct {
	pool db.DB
}

func NewSessionRepository(pool db.DB) workflow.SessionRepository {
	return &SessionRepository{pool: pool}
}

func (r *SessionRepository) Create(ctx context.Context, session *workflow.Session, groupID string, workflowID string) error {
	query := `
		INSERT INTO sessions (id, group_id, workflow_id, status, started_at)
		VALUES ($1, $2, $3, $4, NOW())
	`
	// workflowID can be null if it's a dynamic one not yet saved, but usually it's saved.
	var wfID interface{} = workflowID
	if workflowID == "" {
		wfID = nil
	}

	_, err := r.pool.Exec(ctx, query, session.ID, groupID, wfID, string(session.Status))
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

func (r *SessionRepository) Get(ctx context.Context, id string) (*workflow.SessionEntity, error) {
	query := `
		SELECT id, group_id, workflow_id, status FROM sessions WHERE id = $1
	`
	var s workflow.SessionEntity
	var wfID interface{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&s.ID, &s.GroupID, &wfID, &s.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	if wfID != nil {
		s.WorkflowID = fmt.Sprintf("%v", wfID)
	}
	return &s, nil
}

func (r *SessionRepository) UpdateStatus(ctx context.Context, id string, status workflow.SessionStatus) error {
	query := `
		UPDATE sessions SET status = $2 WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query, id, string(status))
	return err
}
