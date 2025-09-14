package config

import (
	"os"
)

var (
	ApiKeyCoinMarketCap   string
	UrlCoinMarketCap      string
	RedisAddr             string
	ElasticsearchUrl      string
	ElasticsearchUsername string
	ElasticsearchPassword string
)

func LoadEnvs() {
	ApiKeyCoinMarketCap = getEnv("API_KEY_COIN_MARKET_CAP", "")
	UrlCoinMarketCap = getEnv("URL_COIN_MARKET_CAP", "")
	RedisAddr = getEnv("REDIS_ADDR", "localhost:6379")
	ElasticsearchUrl = getEnv("ELASTIC_SEARCH_URL", "http://localhost:9200")
	ElasticsearchUsername = getEnv("ELASTIC_SEARCH_USERNAME", "")
	ElasticsearchPassword = getEnv("ELASTIC_SEARCH_PASSWORD", "")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
