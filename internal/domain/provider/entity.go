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
