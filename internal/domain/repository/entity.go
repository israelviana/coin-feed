package repository

import "time"

type CryptoCurrencyData struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Symbol      string    `json:"symbol"`
	Open        float64   `json:"open"`
	High        float64   `json:"high"`
	Low         float64   `json:"low"`
	Close       float64   `json:"close"`
	Volume      int64     `json:"volume"`
	LastUpdated time.Time `json:"last_updated"`
}
