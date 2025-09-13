package usecase

import (
	"context"
	"strconv"
	"sync"
	"sync/atomic"

	"coin-feed/internal/domain/provider"
	"coin-feed/internal/domain/repository"
)

type iSaveLatestCryptoCurrencyProvider interface {
	FetchLatestCryptoCurrency(ctx context.Context, ids []string) (*provider.LatestCryptoCurrencyResponse, error)
}

type iSaveLatestCryptoCurrencyRepository interface {
	SaveLatestCryptoCurrency(ctx context.Context)
}

type SaveLatestCryptoCurrency struct {
	providerDomain iSaveLatestCryptoCurrencyProvider
	//repository     repository.IRepository
	cache repository.ICacheRepository
}

func NewSaveLatestCryptoCurrency(pd iSaveLatestCryptoCurrencyProvider, ch repository.ICacheRepository) *SaveLatestCryptoCurrency {
	return &SaveLatestCryptoCurrency{
		providerDomain: pd,
		cache:          ch,
	}
}

func (uc *SaveLatestCryptoCurrency) Run(ctx context.Context) error {
	const cacheKey = "cryptoCurrencyMap"
	var response provider.CryptoCurrencyMapResponse
	err := uc.cache.Get(ctx, cacheKey, &response)
	if err != nil {
		return err
	}

	n := len(response.Data)
	if n > 1000 {
		n = 99
	}

	top := response.Data[:n]

	ids := make([]string, n)
	for i, c := range top {
		ids[i] = strconv.Itoa(c.Id)
	}

	latestCryptoCurrency, err := uc.providerDomain.FetchLatestCryptoCurrency(ctx, ids)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	var next int64
	cryptoCurrencyData := make([]*repository.CryptoCurrencyData, len(latestCryptoCurrency.Data))
	for _, cryptoData := range latestCryptoCurrency.Data {
		cryptoData := cryptoData

		wg.Add(1)
		go func() {
			defer wg.Done()
			i := int(atomic.AddInt64(&next, 1)) - 1

			cryptoCurrencyData[i] = &repository.CryptoCurrencyData{
				Id:          cryptoData.Id,
				Name:        cryptoData.Name,
				Symbol:      cryptoData.Symbol,
				Open:        cryptoData.Quote.USD.Open,
				High:        cryptoData.Quote.USD.High,
				Low:         cryptoData.Quote.USD.Low,
				Close:       cryptoData.Quote.USD.Close,
				Volume:      cryptoData.Quote.USD.Volume,
				LastUpdated: cryptoData.LastUpdated,
			}
		}()
	}
	wg.Wait()

	const cacheKeyElastic = "elasticSearch"
	err = uc.cache.Set(ctx, cacheKeyElastic, cryptoCurrencyData, 44640)
	if err != nil {
		return err
	}

	return nil
}
