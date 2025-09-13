package repository

import (
	"context"
)

type ICacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, ttlSeconds int) error
	Get(ctx context.Context, key string, result interface{}) error
	Del(ctx context.Context, key string) error
}

type IRepository interface {
	SaveLatestCryptoCurrency(ctx context.Context)
}
