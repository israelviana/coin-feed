package usecase

import (
	"context"
	"fmt"

	"coin-feed/internal/domain/repository"
	"coin-feed/pkg/logger"

	"go.uber.org/zap"
)

type iGetLatestCryptoCurrencyDataByIdRepository interface {
	GetLatestCryptoCurrencyDataById(ctx context.Context, id string) (*repository.CryptoCurrencyData, error)
}

type GetLatestCryptoCurrencyDataById struct {
	repository iGetLatestCryptoCurrencyDataByIdRepository
	cache      repository.ICacheRepository
}

func NewGetLatestCryptoCurrencyDataById(pd iGetLatestCryptoCurrencyDataByIdRepository, cache repository.ICacheRepository) *GetLatestCryptoCurrencyDataById {
	return &GetLatestCryptoCurrencyDataById{
		repository: pd,
		cache:      cache,
	}
}

func (uc *GetLatestCryptoCurrencyDataById) Run(ctx context.Context, id string) (*repository.CryptoCurrencyData, error) {
	var cacheKey = fmt.Sprintf("cryptoCurrencyData-%s", id)

	var response repository.CryptoCurrencyData
	err := uc.cache.Get(ctx, cacheKey, &response)
	if err != nil {
		cmr, err := uc.repository.GetLatestCryptoCurrencyDataById(ctx, id)
		if err != nil {
			return nil, err
		}

		err = uc.cache.Set(ctx, cacheKey, cmr, 21600)
		if err != nil {
			logger.Logger.Error("Failed to set cache", zap.Error(err))
			return nil, err
		}
		logger.Logger.Info("Succeed to set cache", zap.Error(err))

		return cmr, nil
	}

	return &response, nil
}
