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
	pool     db.DB
	cache    cache.Cache
}

func NewService(embedder llm.Embedder, pool db.DB, cache cache.Cache) *Service {
	return &Service{
		Embedder: embedder,
		pool:     pool,
		cache:    cache,
	}
}

// LogQuarantine writes to PostgreSQL quarantine_logs table (Tier 1)
func (s *Service) LogQuarantine(ctx context.Context, sessionID string, nodeID string, content string, metadata map[string]interface{}) error {
	if s.pool == nil {
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

	_, err = s.pool.Exec(ctx, query, sessionID, content, metaJSON)
	if err != nil {
		return fmt.Errorf("failed to insert quarantine log: %w", err)
	}
	return nil
}

// UpdateWorkingMemory writes to Redis with Ingress Filter (Tier 2)
func (s *Service) UpdateWorkingMemory(ctx context.Context, groupID string, content string, metadata map[string]interface{}) error {
	if s.cache == nil {
		return fmt.Errorf("redis client not initialized")
	}

	// 1. Ingress Filter: Check confidence score
	confidence, ok := metadata["confidence"].(float64)
	if ok && confidence < 0.8 {
		// Reject low-confidence content
		return nil
	}

	// 2. Ingress Filter: LLM-based consistency check (optional, if embedder is available)
	if s.Embedder != nil {
		// Simple heuristic: too short content likely noise
		if len(content) < 50 {
			return nil // Skip very short content
		}
		// In production, we'd call LLM to verify consistency
		// For MVP, we pass if content length is reasonable
	}

	// 3. Write to Redis List
	key := fmt.Sprintf("wm:%s", groupID)

	if err := s.cache.LPush(ctx, key, content).Err(); err != nil {
		return fmt.Errorf("failed to push to working memory: %w", err)
	}

	// Set TTL (24h)
	s.cache.Expire(ctx, key, 24*time.Hour)

	// Cap list size: keep last 50 items
	s.cache.LTrim(ctx, key, 0, 49)

	return nil
}

// CleanupWorkingMemory removes expired entries (called by scheduler)
func (s *Service) CleanupWorkingMemory(ctx context.Context) error {
	if s.cache == nil {
		return fmt.Errorf("redis client not initialized")
	}

	// Redis handles TTL automatically, but we can prune manually if needed
	// This is a placeholder for any additional cleanup logic
	// For example, scanning and removing very old keys:
	// iter := client.Scan(ctx, 0, "wm:*", 100).Iterator()
	// Redis' native TTL handles expiration, so this function primarily exists
	// for potential future enhancements (e.g., archiving before deletion)

	return nil
}

func (s *Service) Promote(ctx context.Context, groupID string, content string) error {
	// 1. Split Text
	splitter := NewRecursiveCharacterSplitter(500, 50) // 500 chars chunk, 50 overlap
	chunks := splitter.SplitText(content)

	if len(chunks) == 0 {
		return nil
	}

	if s.pool == nil {
		return fmt.Errorf("database pool not initialized")
	}

	// 2. Embed and Store Loop
	// Optimization: Batch embedding if provider supports it, but for now loop is simpler for MVP
	for _, chunk := range chunks {
		embedding, err := s.Embedder.Embed(ctx, "default", chunk) // Use default model from config (implied)
		if err != nil {
			return fmt.Errorf("failed to embed chunk: %w", err)
		}

		vecBytes, _ := json.Marshal(embedding)
		vecStr := string(vecBytes)

		query := `
			INSERT INTO memories (group_id, content, embedding, metadata)
			VALUES ($1, $2, $3, $4)
		`
		// Metadata can store source info
		meta := map[string]interface{}{
			"source":      "promotion",
			"promoted_at": time.Now(),
		}
		metaJSON, _ := json.Marshal(meta)

		_, err = s.pool.Exec(ctx, query, groupID, chunk, vecStr, metaJSON)
		if err != nil {
			return fmt.Errorf("failed to store memory chunk: %w", err)
		}
	}

	return nil
}

func (s *Service) Retrieve(ctx context.Context, query string, groupID string, sessionID string) ([]ContextItem, error) {
	var items []ContextItem

	// 1. Hot Working Memory (Redis)
	if s.cache != nil {
		key := fmt.Sprintf("wm:%s", groupID)
		// Get last 10 items
		vals, err := s.cache.LRange(ctx, key, 0, 10).Result()
		if err == nil {
			for _, v := range vals {
				items = append(items, ContextItem{Content: v, Source: "hot", Score: 1.0})
			}
		}
	}

	// 2. Cold LTM (PGVector)
	if s.pool != nil && s.Embedder != nil {
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
			q := `SELECT content, 1 - (embedding <=> $1) as score FROM memories WHERE group_id = $2::uuid`
			params := []interface{}{vecStr, groupID}
			if sessionID != "" {
				q += ` AND session_id = $3::uuid`
				params = append(params, sessionID)
			}
			q += ` ORDER BY embedding <=> $1 LIMIT 5`

			rows, err := s.pool.Query(ctx, q, params...)
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
