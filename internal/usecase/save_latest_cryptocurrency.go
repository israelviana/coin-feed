package usecase

import (
	"context"
	"sync"
	"time"

	"coin-feed/internal/domain/provider"
	"coin-feed/internal/domain/repository"
)

type iSaveLatestCryptoCurrencyProvider interface {
	FetchLatestCryptoCurrency(ctx context.Context) (*provider.LatestCryptoCurrencyResponse, error)
}

type iSaveLatestCryptoCurrencyRepository interface {
	SaveLatestCryptoCurrency(ctx context.Context, data []*repository.CryptoCurrencyData) error
}

type SaveLatestCryptoCurrency struct {
	providerDomain iSaveLatestCryptoCurrencyProvider
	repository     iSaveLatestCryptoCurrencyRepository
}

func NewSaveLatestCryptoCurrency(pd iSaveLatestCryptoCurrencyProvider, rp iSaveLatestCryptoCurrencyRepository) *SaveLatestCryptoCurrency {
	return &SaveLatestCryptoCurrency{
		providerDomain: pd,
		repository:     rp,
	}
}

func (uc *SaveLatestCryptoCurrency) Run(ctx context.Context) error {
	defer ctx.Done()

	latestCryptoCurrency, err := uc.providerDomain.FetchLatestCryptoCurrency(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	cryptoCurrencyData := make([]*repository.CryptoCurrencyData, len(latestCryptoCurrency.Data))
	for i, cryptoData := range latestCryptoCurrency.Data {
		cryptoData := cryptoData
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()

			cryptoCurrencyData[i] = &repository.CryptoCurrencyData{
				Id:                cryptoData.Id,
				Name:              cryptoData.Name,
				Symbol:            cryptoData.Symbol,
				Platform:          cryptoData.Platform,
				TotalSupply:       cryptoData.TotalSupply,
				CirculatingSupply: cryptoData.CirculatingSupply,
				Price:             cryptoData.Quote.USD.Price,
				Volume24H:         cryptoData.Quote.USD.Volume24H,
				VolumeChange24H:   cryptoData.Quote.USD.VolumeChange24H,
				PercentChange1H:   cryptoData.Quote.USD.PercentChange1H,
				PercentChange24H:  cryptoData.Quote.USD.PercentChange24H,
				PercentChange7D:   cryptoData.Quote.USD.PercentChange7D,
				PercentChange30D:  cryptoData.Quote.USD.PercentChange30D,
				PercentChange60D:  cryptoData.Quote.USD.PercentChange60D,
				PercentChange90D:  cryptoData.Quote.USD.PercentChange90D,
				LastUpdated:       cryptoData.LastUpdated,
				CreatedAt:         time.Now(),
			}
		}()
	}
	wg.Wait()

	err = uc.repository.SaveLatestCryptoCurrency(ctx, cryptoCurrencyData)
	if err != nil {
		return err
	}

	return nil
}
