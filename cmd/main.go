package main

import (
	"context"

	"coin-feed/cmd/api"
	"coin-feed/cmd/job"
	"coin-feed/config"
	"coin-feed/internal/usecase"
	pkgElasticsearch "coin-feed/pkg/elasticsearch"
	"coin-feed/pkg/logger"
	pkgRedis "coin-feed/pkg/redis"
	cmcProvider "coin-feed/providers/coinmarketcap"
	elasticsearchRepo "coin-feed/repositories/elasticsearch"
	redisRepo "coin-feed/repositories/redis"

	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()
	defer logger.Logger.Sync()

	config.LoadEnvs()
	coinMarketCapProvider := cmcProvider.NewProvider(config.UrlCoinMarketCap, config.ApiKeyCoinMarketCap)
	redisClient := pkgRedis.NewRedisClient()
	redisRepository := redisRepo.NewRedisRepository(redisClient)
	elasticsearchClient, err := pkgElasticsearch.NewClient([]string{config.ElasticsearchUrl}, config.ElasticsearchUsername, config.ElasticsearchPassword, nil)
	if err != nil {
		logger.Logger.Fatal("Failed to create elasticsearch client", zap.Error(err))
	}
	elasticsearchRepository := elasticsearchRepo.NewRepository(elasticsearchClient, "tito-crypto-tracker")
	fetchCryptoCurrencyMapUseCase := usecase.NewFetchCryptocurrencyMap(coinMarketCapProvider, redisRepository)
	saveLatestCryptoUseCase := usecase.NewSaveLatestCryptoCurrency(coinMarketCapProvider, elasticsearchRepository)
	cryptoHandler := api.NewCryptoHandler(fetchCryptoCurrencyMapUseCase)

	//save crypto data
	go func() {
		job.Start(context.Background(), saveLatestCryptoUseCase)
	}()

	//gin api
	api.Start(cryptoHandler)
}
