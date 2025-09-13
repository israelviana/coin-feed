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
	Data   map[string]Data `json:"data"`
	Status Status          `json:"status"`
}

type Data struct {
	Id          int         `json:"id"`
	Name        string      `json:"name"`
	Symbol      string      `json:"symbol"`
	LastUpdated time.Time   `json:"last_updated"`
	TimeOpen    time.Time   `json:"time_open"`
	TimeClose   interface{} `json:"time_close"`
	TimeHigh    time.Time   `json:"time_high"`
	TimeLow     time.Time   `json:"time_low"`
	Quote       Quote       `json:"quote"`
}

type Quote struct {
	USD USD `json:"USD"`
}

type USD struct {
	Open        float64   `json:"open"`
	High        float64   `json:"high"`
	Low         float64   `json:"low"`
	Close       float64   `json:"close"`
	Volume      int64     `json:"volume"`
	LastUpdated time.Time `json:"last_updated"`
}

func (l *LatestCryptoCurrencyResponse) MarshalBinary() ([]byte, error) {
	return json.Marshal(l)
}

func (l *LatestCryptoCurrencyResponse) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, l)
}
