package memory

import (
	"context"
	"testing"

	"github.com/hrygo/council/internal/infrastructure/cache"
	"github.com/hrygo/council/internal/infrastructure/llm"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/redis/go-redis/v9"
)

func TestService_LogQuarantine(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()

	svc := NewService(&llm.MockProvider{}, mockDB, nil)

	sessionID := "session-1"
	nodeID := "node-1"
	content := "harmful content"
	metadata := map[string]interface{}{"reason": "safety"}

	mockDB.ExpectExec("INSERT INTO quarantine_logs").
		WithArgs(sessionID, content, pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = svc.LogQuarantine(context.Background(), sessionID, nodeID, content, metadata)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestService_UpdateWorkingMemory(t *testing.T) {
	mockCache := &cache.MockCache{}
	svc := NewService(nil, nil, mockCache)

	groupID := "group-1"
	content := "This is a long enough content to pass the ingress filter of 50 chars minimum."
	metadata := map[string]interface{}{"confidence": 0.9}

	calledLPush := false
	mockCache.LPushFunc = func(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
		calledLPush = true
		if key != "wm:group-1" {
			t.Errorf("expected key wm:group-1, got %s", key)
		}
		return redis.NewIntCmd(ctx)
	}

	err := svc.UpdateWorkingMemory(context.Background(), groupID, content, metadata)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if !calledLPush {
		t.Error("expected LPush to be called")
	}
}

func TestService_Promote(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	mockEmbedder := &llm.MockProvider{
		EmbedResponse: []float32{0.1, 0.2, 0.3},
	}
	svc := NewService(mockEmbedder, mockDB, nil)

	content := "Short content for single chunk"
	groupID := "group-1"

	mockDB.ExpectExec("INSERT INTO memories").
		WithArgs(groupID, content, pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := svc.Promote(context.Background(), groupID, content)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestService_Retrieve(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	mockCache := &cache.MockCache{}
	mockEmbedder := &llm.MockProvider{
		EmbedResponse: []float32{0.1, 0.2, 0.3},
	}
	svc := NewService(mockEmbedder, mockDB, mockCache)

	groupID := "group-1"
	query := "test query"

	// Mock Cache LRange
	mockCache.LRangeFunc = func(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
		cmd := redis.NewStringSliceCmd(ctx)
		cmd.SetVal([]string{"hot data"})
		return cmd
	}

	// Mock DB Query for PGVector
	mockDB.ExpectQuery("SELECT content, 1 - \\(embedding <=> \\$1\\) as score FROM memories").
		WithArgs(pgxmock.AnyArg(), groupID).
		WillReturnRows(pgxmock.NewRows([]string{"content", "score"}).AddRow("cold data", 0.95))

	items, err := svc.Retrieve(context.Background(), query, groupID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
	if items[0].Source != "hot" || items[1].Source != "cold" {
		t.Errorf("expected hot then cold, got %s, %s", items[0].Source, items[1].Source)
	}
}
