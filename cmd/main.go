package main

import (
	"context"

	"coin-feed/cmd/api"
	"coin-feed/cmd/job"
	"coin-feed/config"
	"coin-feed/internal/usecase"
	pkgRedis "coin-feed/pkg/redis"
	cmcProvider "coin-feed/providers/coinmarketcap"
	redisRepo "coin-feed/repositories/redis"
)

func main() {
	coinMarketCapProvider := cmcProvider.NewProvider(config.UrlCoinMarketCap, config.ApiKeyCoinMarketCap)
	redisClient := pkgRedis.NewRedisClient()
	redisRepository := redisRepo.NewRedisRepository(redisClient)
	fetchCryptoCurrencyMapUseCase := usecase.NewFetchCryptocurrencyMap(coinMarketCapProvider, redisRepository)
	saveLatestCryptoUseCase := usecase.NewSaveLatestCryptoCurrency(coinMarketCapProvider, redisRepository)
	cryptoHandler := api.NewCryptoHandler(fetchCryptoCurrencyMapUseCase)

	//save crypto data
	go func() {
		job.Start(context.Background(), saveLatestCryptoUseCase)
	}()

	//gin api
	api.Start(cryptoHandler)
}
