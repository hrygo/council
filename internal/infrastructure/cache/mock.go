package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// MockCache implements Cache interface for testing.
type MockCache struct {
	LPushFunc  func(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	LRangeFunc func(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd
	LTrimFunc  func(ctx context.Context, key string, start, stop int64) *redis.StatusCmd
	ExpireFunc func(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
	DelFunc    func(ctx context.Context, keys ...string) *redis.IntCmd
}

func (m *MockCache) LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	if m.LPushFunc != nil {
		return m.LPushFunc(ctx, key, values...)
	}
	return redis.NewIntCmd(ctx)
}

func (m *MockCache) LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	if m.LRangeFunc != nil {
		return m.LRangeFunc(ctx, key, start, stop)
	}
	res := redis.NewStringSliceCmd(ctx)
	res.SetVal([]string{})
	return res
}

func (m *MockCache) LTrim(ctx context.Context, key string, start, stop int64) *redis.StatusCmd {
	if m.LTrimFunc != nil {
		return m.LTrimFunc(ctx, key, start, stop)
	}
	return redis.NewStatusCmd(ctx)
}

func (m *MockCache) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	if m.ExpireFunc != nil {
		return m.ExpireFunc(ctx, key, expiration)
	}
	return redis.NewBoolCmd(ctx)
}

func (m *MockCache) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	if m.DelFunc != nil {
		return m.DelFunc(ctx, keys...)
	}
	return redis.NewIntCmd(ctx)
}
