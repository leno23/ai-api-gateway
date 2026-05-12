package store

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func OpenRedis(addr string) (*redis.Client, error) {
	c := redis.NewClient(&redis.Options{Addr: addr})
	if err := c.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}
	return c, nil
}
