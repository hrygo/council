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

// Init initializes the Redis client (singleton)
func Init(addr string, password string, db int) error {
	var err error
	once.Do(func() {
		client, err = connect(addr, password, db)
	})
	return err
}

func connect(addr string, password string, db int) (*redis.Client, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if pingErr := c.Ping(context.Background()).Err(); pingErr != nil {
		c.Close()
		return nil, fmt.Errorf("failed to connect to redis: %w", pingErr)
	}

	return c, nil
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
