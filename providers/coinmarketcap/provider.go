package coinmarketcap

import (
	"context"
	"strings"
	"time"

	"coin-feed/internal/domain/provider"

	"github.com/go-resty/resty/v2"
)

type Provider struct {
	resty *resty.Client
}

func NewProvider(url, apiKey string) *Provider {
	return &Provider{
		resty: resty.New().
			SetBaseURL(url).
			SetHeader("X-CMC_PRO_API_KEY", apiKey).
			SetHeader("Accept", "application/json").
			SetTimeout(15 * time.Second).
			SetRetryCount(3).
			SetRetryWaitTime(500 * time.Millisecond).
			SetRetryMaxWaitTime(5 * time.Second),
	}
}

func (c *Provider) FetchCryptoCurrencyMap(ctx context.Context) (*provider.CryptoCurrencyMapResponse, error) {
	var cryptoCurrencyMap provider.CryptoCurrencyMapResponse
	_, err := c.resty.R().SetContext(ctx).SetResult(&cryptoCurrencyMap).
		SetQueryParams(map[string]string{
			"sort": "cmc_rank",
		}).
		Get("/v1/cryptocurrency/map")
	if err != nil {
		return nil, err
	}

	return &cryptoCurrencyMap, nil
}

func (c *Provider) FetchLatestCryptoCurrency(ctx context.Context, ids []string) (*provider.LatestCryptoCurrencyResponse, error) {
	var cryptoCurrencyMap provider.LatestCryptoCurrencyResponse
	_, err := c.resty.R().SetContext(ctx).SetResult(&cryptoCurrencyMap).
		SetQueryParams(map[string]string{
			"id": strings.Join(ids, ","),
		}).
		Get("/v2/cryptocurrency/ohlcv/latest")
	if err != nil {
		return nil, err
	}

	return &cryptoCurrencyMap, nil
}
