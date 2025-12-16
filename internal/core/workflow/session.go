package workflow

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type SessionStatus string

const (
	SessionPending   SessionStatus = "pending"
	SessionRunning   SessionStatus = "running"
	SessionPaused    SessionStatus = "paused"
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

	ctx      context.Context
	cancel   context.CancelFunc
	resumeCh chan struct{}

	SignalChannels map[string]chan interface{}
	mu             sync.RWMutex
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
	s.resumeCh = make(chan struct{})
	close(s.resumeCh) // Initially not paused

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

func (s *Session) Pause() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Status == SessionRunning {
		s.Status = SessionPaused
		s.resumeCh = make(chan struct{}) // Create a new blocking channel
	}
}

func (s *Session) Resume() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Status == SessionPaused {
		s.Status = SessionRunning
		close(s.resumeCh) // Unblock all waiters
	}
}

func (s *Session) WaitIfPaused(ctx context.Context) error {
	s.mu.RLock()
	ch := s.resumeCh
	status := s.Status
	s.mu.RUnlock()

	if status != SessionPaused {
		return nil
	}

	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Signal support
func (s *Session) GetSignalChannel(nodeID string) chan interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.SignalChannels == nil {
		s.SignalChannels = make(map[string]chan interface{})
	}
	if _, ok := s.SignalChannels[nodeID]; !ok {
		s.SignalChannels[nodeID] = make(chan interface{}, 1)
	}
	return s.SignalChannels[nodeID]
}

func (s *Session) SendSignal(nodeID string, payload interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.SignalChannels == nil {
		return fmt.Errorf("no signal channels active")
	}
	ch, ok := s.SignalChannels[nodeID]
	if !ok {
		return fmt.Errorf("node %s is not waiting for signal", nodeID)
	}

	// Non-blocking send or blocking? Ideally blocking if buffer is full,
	// but here we use buffered channel of 1 usually.
	select {
	case ch <- payload:
		return nil
	default:
		return fmt.Errorf("signal channel full or receiver not ready")
	}
}

func (s *Session) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Allow stopping from Paused state too
	if s.Status == SessionRunning || s.Status == SessionPaused {
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
