package redis

import (
	"context"

	"coin-feed/config"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	addr := config.RedisAddr
	db := 0

	return redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})
}

func Ping(ctx context.Context, client *redis.Client) error {
	_, err := client.Ping(ctx).Result()
	return err
}
