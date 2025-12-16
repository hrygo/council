package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hrygo/council/internal/infrastructure/cache"
	"github.com/hrygo/council/internal/infrastructure/db"
	"github.com/hrygo/council/internal/infrastructure/llm"
)

type Service struct {
	Embedder llm.Embedder
}

func NewService(embedder llm.Embedder) *Service {
	return &Service{Embedder: embedder}
}

// LogQuarantine writes to PostgreSQL quarantine_logs table (Tier 1)
func (s *Service) LogQuarantine(ctx context.Context, sessionID string, nodeID string, content string, metadata map[string]interface{}) error {
	pool := db.GetPool()
	if pool == nil {
		return fmt.Errorf("database pool not initialized")
	}

	metaJSON, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `
		INSERT INTO quarantine_logs (session_id, content, raw_metadata)
		VALUES ($1, $2, $3)
	`
	// Note: We might want to include node_id in metadata or table.
	// Current schema has specific fields. Let's stick to schema in migration.
	// Schema: id, session_id, content, raw_metadata.

	// We inject node_id into metadata for storage
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	metadata["node_id"] = nodeID

	_, err = pool.Exec(ctx, query, sessionID, content, metaJSON)
	if err != nil {
		return fmt.Errorf("failed to insert quarantine log: %w", err)
	}
	return nil
}

// UpdateWorkingMemory writes to Redis with Ingress Filter (Tier 2)
func (s *Service) UpdateWorkingMemory(ctx context.Context, groupID string, content string, metadata map[string]interface{}) error {
	client := cache.GetClient()
	if client == nil {
		return fmt.Errorf("redis client not initialized")
	}

	// 1. Ingress Filter
	// Check confidence score
	confidence, ok := metadata["confidence"].(float64)
	if ok && confidence < 0.8 {
		// Reject (Log rejection?)
		return nil
	}

	// 2. Write to Redis List or Set
	// Key: wm:{group_id}
	// Value: content (or JSON with metadata)
	key := fmt.Sprintf("wm:%s", groupID)

	// For MVP: Just push to a list, TTL handled by expiration policy or manually?
	// Redis List approach: LPUSH
	if err := client.LPush(ctx, key, content).Err(); err != nil {
		return fmt.Errorf("failed to push to working memory: %w", err)
	}

	// Set TTL if new key (24h)
	client.Expire(ctx, key, 24*time.Hour)

	// Cap list size? keeping last 50 items
	client.LTrim(ctx, key, 0, 49)

	return nil
}

func (s *Service) Promote(ctx context.Context, groupID string, digest string) error {
	// Stub for Tier 3
	return nil
}

func (s *Service) Retrieve(ctx context.Context, query string, groupID string) ([]ContextItem, error) {
	var items []ContextItem

	// 1. Hot Working Memory (Redis)
	if client := cache.GetClient(); client != nil {
		key := fmt.Sprintf("wm:%s", groupID)
		// Get last 10 items
		vals, err := client.LRange(ctx, key, 0, 10).Result()
		if err == nil {
			for _, v := range vals {
				items = append(items, ContextItem{Content: v, Source: "hot", Score: 1.0})
			}
		}
	}

	// 2. Cold LTM (PGVector)
	if pool := db.GetPool(); pool != nil && s.Embedder != nil {
		// Generate Embedding
		// Default model: text-embedding-ada-002 or compatible
		embedding, err := s.Embedder.Embed(ctx, "text-embedding-ada-002", query)
		if err == nil {
			// Convert []float32 to vector string "[0.1,0.2,...]"
			var vecStr string
			// optimization: build string manually or json marshal?
			// PGVector expects '[...]'
			// We construct it simply
			bytes, _ := json.Marshal(embedding)
			vecStr = string(bytes)

			// Query
			// Using <-> (L2 distance) or <=> (Cosine distance).
			// Cosine distance is 1 - Cosine Similarity.
			q := `SELECT content, 1 - (embedding <=> $1) as score FROM memories WHERE group_id = $2::uuid ORDER BY embedding <=> $1 LIMIT 5`

			rows, err := pool.Query(ctx, q, vecStr, groupID)
			if err == nil {
				defer rows.Close()
				for rows.Next() {
					var content string
					var score float64
					if err := rows.Scan(&content, &score); err == nil {
						items = append(items, ContextItem{Content: content, Source: "cold", Score: score})
					}
				}
			}
		}
	}

	return items, nil
}
