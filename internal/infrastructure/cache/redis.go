package cache

import (
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	once   sync.Once
)

// Init initializes the Redis client
func Init(addr string, password string, db int) error {
	var err error
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		})

		if pingErr := client.Ping(context.Background()).Err(); pingErr != nil {
			err = fmt.Errorf("failed to connect to redis: %w", pingErr)
		}
	})
	return err
}

// GetClient returns the Redis client
func GetClient() *redis.Client {
	return client
}

// Close closes the Redis connection
func Close() {
	if client != nil {
		client.Close()
	}
}
