package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache defines the caching interface for services.
type Cache interface {
	LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd
	LTrim(ctx context.Context, key string, start, stop int64) *redis.StatusCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}
