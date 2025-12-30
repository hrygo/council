package persistence

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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
		INSERT INTO sessions (session_uuid, group_uuid, workflow_uuid, status, started_at)
		VALUES ($1, $2, $3, $4, NOW())
	`
	// Handle empty or invalid UUIDs - PostgreSQL expects NULL for empty/invalid UUID values
	var grpID interface{} = nil
	if groupID != "" && isValidUUID(groupID) {
		grpID = groupID
	}
	var wfID interface{} = nil
	if workflowID != "" && isValidUUID(workflowID) {
		wfID = workflowID
	}

	_, err := r.pool.Exec(ctx, query, session.ID, grpID, wfID, string(session.Status))
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

func (r *SessionRepository) Get(ctx context.Context, id string) (*workflow.SessionEntity, error) {
	query := `
		SELECT session_uuid, group_uuid, workflow_uuid, status FROM sessions WHERE session_uuid = $1
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
		UPDATE sessions SET status = $2 WHERE session_uuid = $1
	`
	_, err := r.pool.Exec(ctx, query, id, string(status))
	return err
}

// isValidUUID checks if a string is a valid UUID format
func isValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}
