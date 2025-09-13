package provider

import (
	"context"
)

type IProvider interface {
	FetchCryptoCurrencyMap(ctx context.Context) (*CryptoCurrencyMapResponse, error)
	FetchLatestCryptoCurrency(ctx context.Context) (*LatestCryptoCurrencyResponse, error)
}
