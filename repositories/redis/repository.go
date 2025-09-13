package redis

import (
	"context"
	"encoding/json"
	"time"

	"coin-feed/internal/domain/repository"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) repository.ICacheRepository {
	return &RedisRepository{client: client}
}

func (r *RedisRepository) Set(ctx context.Context, key string, value interface{}, ttlSeconds int) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, jsonData, time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *RedisRepository) Get(ctx context.Context, key string, result interface{}) error {
	return r.client.Get(ctx, key).Scan(result)
}

func (r *RedisRepository) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
