package provider

import (
	"context"
)

type IProviderCrypto interface {
	FetchCryptoCurrencyMap(ctx context.Context) (*CryptoCurrencyMapResponse, error)
}
