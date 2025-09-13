package usecase

import (
	"context"

	"coin-feed/internal/domain/provider"
	"coin-feed/internal/domain/repository"
	"coin-feed/pkg/logger"

	"go.uber.org/zap"
)

type FetchCryptocurrencyMap struct {
	providerDomain provider.IProviderCrypto
	redis          repository.ICacheRepository
}

func NewFetchCryptocurrencyMap(pd provider.IProviderCrypto, cache repository.ICacheRepository) *FetchCryptocurrencyMap {
	return &FetchCryptocurrencyMap{
		providerDomain: pd,
		redis:          cache,
	}
}

func (uc *FetchCryptocurrencyMap) Fetch(ctx context.Context) (*provider.CryptoCurrencyMapResponse, error) {
	const cacheKey = "cryptoCurrencyMap"

	var response provider.CryptoCurrencyMapResponse
	err := uc.redis.Get(ctx, cacheKey, &response)
	if err != nil {
		cmr, err := uc.providerDomain.FetchCryptoCurrencyMap(ctx)
		if err != nil {
			return nil, err
		}

		err = uc.redis.Set(ctx, cacheKey, cmr, 44640)
		if err != nil {
			logger.Logger.Error("Failed to set cache", zap.Error(err))
			return nil, err
		}
		logger.Logger.Info("Succeed to set cache", zap.Error(err))

		return cmr, nil
	}

	return &response, nil
}
