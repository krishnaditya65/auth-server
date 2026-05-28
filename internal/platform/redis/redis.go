package redis

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"
)

func New(addr string) (*goredis.Client, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr: addr,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return client, nil
}
