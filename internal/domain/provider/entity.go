package provider

import (
	"encoding/json"
	"time"
)

type Status struct {
	Timestamp    time.Time   `json:"timestamp"`
	ErrorCode    int         `json:"error_code"`
	ErrorMessage interface{} `json:"error_message"`
	Elapsed      int         `json:"elapsed"`
	CreditCount  int         `json:"credit_count"`
	Notice       interface{} `json:"notice"`
}

type CryptoCurrencyMapResponse struct {
	Status Status           `json:"status"`
	Data   []CryptoCurrency `json:"data"`
}

type CryptoCurrency struct {
	Id                  int       `json:"id"`
	Rank                int       `json:"rank"`
	Name                string    `json:"name"`
	Symbol              string    `json:"symbol"`
	Slug                string    `json:"slug"`
	IsActive            int       `json:"is_active"`
	Status              int       `json:"status"`
	FirstHistoricalData time.Time `json:"first_historical_data"`
	LastHistoricalData  time.Time `json:"last_historical_data"`
	Platform            Platform  `json:"platform"`
}

type Platform struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	Slug         string `json:"slug"`
	TokenAddress string `json:"token_address"`
}

func (c *CryptoCurrencyMapResponse) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c *CryptoCurrencyMapResponse) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, c)
}

type LatestCryptoCurrencyResponse struct {
	Data   []Data `json:"data"`
	Status Status `json:"status"`
}

type Data struct {
	Id                int         `json:"id"`
	Name              string      `json:"name"`
	Symbol            string      `json:"symbol"`
	Platform          interface{} `json:"platform"`
	TotalSupply       float64     `json:"total_supply"`
	CirculatingSupply float64     `json:"circulating_supply"`
	Quote             Quote       `json:"quote"`
	LastUpdated       time.Time   `json:"last_updated"`
}

type Quote struct {
	USD USD `json:"USD"`
}

type USD struct {
	Price                 float64     `json:"price"`
	Volume24H             float64     `json:"volume_24h"`
	VolumeChange24H       float64     `json:"volume_change_24h"`
	PercentChange1H       float64     `json:"percent_change_1h"`
	PercentChange24H      float64     `json:"percent_change_24h"`
	PercentChange7D       float64     `json:"percent_change_7d"`
	PercentChange30D      float64     `json:"percent_change_30d"`
	PercentChange60D      float64     `json:"percent_change_60d"`
	PercentChange90D      float64     `json:"percent_change_90d"`
	MarketCap             float64     `json:"market_cap"`
	MarketCapDominance    float64     `json:"market_cap_dominance"`
	FullyDilutedMarketCap float64     `json:"fully_diluted_market_cap"`
	Tvl                   interface{} `json:"tvl"`
	LastUpdated           time.Time   `json:"last_updated"`
}

func (l *LatestCryptoCurrencyResponse) MarshalBinary() ([]byte, error) {
	return json.Marshal(l)
}

func (l *LatestCryptoCurrencyResponse) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, l)
}
