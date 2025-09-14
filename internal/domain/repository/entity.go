package repository

import (
	"encoding/json"
	"time"
)

type CryptoCurrencyData struct {
	Id                int         `json:"id"`
	Name              string      `json:"name"`
	Symbol            string      `json:"symbol"`
	Platform          interface{} `json:"platform"`
	TotalSupply       float64     `json:"total_supply"`
	CirculatingSupply float64     `json:"circulating_supply"`
	Price             float64     `json:"price"`
	Volume24H         float64     `json:"volume_24h"`
	VolumeChange24H   float64     `json:"volume_change_24h"`
	PercentChange1H   float64     `json:"percent_change_1h"`
	PercentChange24H  float64     `json:"percent_change_24h"`
	PercentChange7D   float64     `json:"percent_change_7d"`
	PercentChange30D  float64     `json:"percent_change_30d"`
	PercentChange60D  float64     `json:"percent_change_60d"`
	PercentChange90D  float64     `json:"percent_change_90d"`
	LastUpdated       time.Time   `json:"last_updated"`
	CreatedAt         time.Time   `json:"created_at"`
}

func (c *CryptoCurrencyData) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c *CryptoCurrencyData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, c)
}
