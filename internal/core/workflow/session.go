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
	ID        string `json:"session_uuid"`
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

	SignalChannels map[string]chan interface{} `json:"signal_channels,omitempty"`
	ContextData    map[string]interface{}      `json:"context_data"` // Runtime context for Loop variables, etc.
	FileRepo       SessionFileRepository       `json:"-"`            // Injected persistence
	mu             sync.RWMutex
}

func (s *Session) SetFileRepository(repo SessionFileRepository) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.FileRepo = repo
}

func (s *Session) WriteFile(path, content, author, reason string) (int, error) {
	s.mu.RLock()
	repo := s.FileRepo
	s.mu.RUnlock()

	if repo == nil {
		return 0, fmt.Errorf("file repository not injected")
	}
	version, err := repo.AddVersion(s.Context(), s.ID, path, content, author, reason)
	if err == nil {
		// Broadcast VFS change
		// Need access to Engine Stream? No, Session doesn't have it.
		// We could send to a dedicated VFS channel if we had one.
		// For now, assume Frontend polls or we enhance Engine later.
		// Actually, we can use a callback if we want.
		_ = 0 // SA9003 suppression: intentional empty branch for future logic
	}
	return version, err
}

func (s *Session) GetLatestFile(path string) (*FileEntity, error) {
	s.mu.RLock()
	repo := s.FileRepo
	s.mu.RUnlock()

	if repo == nil {
		return nil, fmt.Errorf("file repository not injected")
	}
	return repo.GetLatest(s.Context(), s.ID, path)
}

func (s *Session) ListFiles() (map[string]*FileEntity, error) {
	s.mu.RLock()
	repo := s.FileRepo
	s.mu.RUnlock()

	if repo == nil {
		return nil, fmt.Errorf("file repository not injected")
	}
	return repo.ListFiles(s.Context(), s.ID)
}

func (s *Session) SetContext(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ContextData == nil {
		s.ContextData = make(map[string]interface{})
	}
	s.ContextData[key] = value
}

func (s *Session) GetContext(key string) interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.ContextData == nil {
		return nil
	}
	return s.ContextData[key]
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
