package workflow

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type SessionStatus string

const (
	SessionPending   SessionStatus = "pending"
	SessionRunning   SessionStatus = "running"
	SessionCompleted SessionStatus = "completed"
	SessionFailed    SessionStatus = "failed"
	SessionCancelled SessionStatus = "cancelled"
)

// Session represents a single execution instance of a workflow
type Session struct {
	ID        string
	Graph     *GraphDefinition
	Status    SessionStatus
	StartTime time.Time
	EndTime   time.Time
	Inputs    map[string]interface{}
	Outputs   map[string]interface{}
	Error     error

	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.RWMutex
}

func NewSession(graph *GraphDefinition, inputs map[string]interface{}) *Session {
	return &Session{
		ID:     uuid.New().String(),
		Graph:  graph,
		Inputs: inputs,
		Status: SessionPending,
	}
}

func (s *Session) Start(parentCtx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Status = SessionRunning
	s.StartTime = time.Now()

	// Create cancelable context for this session
	s.ctx, s.cancel = context.WithCancel(parentCtx)
}

func (s *Session) Complete() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Status = SessionCompleted
	s.EndTime = time.Now()
	if s.cancel != nil {
		s.cancel() // Cleanup resources
	}
}

func (s *Session) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Status == SessionRunning {
		s.Status = SessionFailed // or Cancelled, using Failed based on test expectation
		s.EndTime = time.Now()
	}
	if s.cancel != nil {
		s.cancel()
	}
}

func (s *Session) Context() context.Context {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.ctx == nil {
		return context.Background()
	}
	return s.ctx
}
