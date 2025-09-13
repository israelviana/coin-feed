package usecase

import (
	"context"

	"coin-feed/internal/domain/provider"
	"coin-feed/internal/domain/repository"
	"coin-feed/pkg/logger"

	"go.uber.org/zap"
)

type iFetchCryptoCurrencyMapProvider interface {
	FetchCryptoCurrencyMap(ctx context.Context) (*provider.CryptoCurrencyMapResponse, error)
}

type FetchCryptocurrencyMap struct {
	providerDomain iFetchCryptoCurrencyMapProvider
	cache          repository.ICacheRepository
}

func NewFetchCryptocurrencyMap(pd iFetchCryptoCurrencyMapProvider, cache repository.ICacheRepository) *FetchCryptocurrencyMap {
	return &FetchCryptocurrencyMap{
		providerDomain: pd,
		cache:          cache,
	}
}

func (uc *FetchCryptocurrencyMap) Run(ctx context.Context) (*provider.CryptoCurrencyMapResponse, error) {
	const cacheKey = "cryptoCurrencyMap"

	var response provider.CryptoCurrencyMapResponse
	err := uc.cache.Get(ctx, cacheKey, &response)
	if err != nil {
		cmr, err := uc.providerDomain.FetchCryptoCurrencyMap(ctx)
		if err != nil {
			return nil, err
		}

		err = uc.cache.Set(ctx, cacheKey, cmr, 44640)
		if err != nil {
			logger.Logger.Error("Failed to set cache", zap.Error(err))
			return nil, err
		}
		logger.Logger.Info("Succeed to set cache", zap.Error(err))

		return cmr, nil
	}

	return &response, nil
}
