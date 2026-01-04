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
		INSERT INTO sessions (session_uuid, group_uuid, workflow_uuid, status, proposal, node_statuses, started_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`
	// Handle empty strings - PostgreSQL expects NULL for empty UUID values
	var grpID interface{} = groupID
	if groupID == "" {
		grpID = nil
	}
	var wfID interface{} = workflowID
	if workflowID == "" {
		wfID = nil
	}

	proposal := session.Inputs["proposal"]
	nodeStatuses := session.NodeStatuses
	if nodeStatuses == nil {
		nodeStatuses = make(map[string]workflow.NodeStatus)
	}

	_, err := r.pool.Exec(ctx, query, session.ID, grpID, wfID, string(session.Status), proposal, nodeStatuses)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

func (r *SessionRepository) Get(ctx context.Context, id string) (*workflow.SessionEntity, error) {
	query := `
		SELECT session_uuid, COALESCE(group_uuid::text, ''), COALESCE(workflow_uuid::text, ''), 
		       status, proposal, COALESCE(node_statuses, '{}'::jsonb), started_at, ended_at 
		FROM sessions WHERE session_uuid = $1
	`
	var s workflow.SessionEntity
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.GroupID, &s.WorkflowID, &s.Status, &s.Proposal, &s.NodeStatuses, &s.StartedAt, &s.EndedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	return &s, nil
}

func (r *SessionRepository) UpdateStatus(ctx context.Context, id string, status workflow.SessionStatus) error {
	query := `UPDATE sessions SET status = $2, updated_at = NOW()`
	if status == workflow.SessionCompleted || status == workflow.SessionFailed || status == workflow.SessionCancelled {
		query += ", ended_at = NOW()"
	}
	query += " WHERE session_uuid = $1"
	_, err := r.pool.Exec(ctx, query, id, string(status))
	return err
}

func (r *SessionRepository) UpdateNodeStatus(ctx context.Context, sessionID string, nodeID string, status workflow.NodeStatus) error {
	// Use jsonb || operator to merge/update the status for the specific node
	query := `
		UPDATE sessions 
		SET node_statuses = COALESCE(node_statuses, '{}'::jsonb) || jsonb_build_object($2::text, $3::text),
		    updated_at = NOW()
		WHERE session_uuid = $1
	`
	_, err := r.pool.Exec(ctx, query, sessionID, nodeID, string(status))
	if err != nil {
		return fmt.Errorf("failed to update node status: %w", err)
	}
	return nil
}
