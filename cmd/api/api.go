package api

import (
	"context"
	"os"

	"coin-feed/config"
	"coin-feed/internal/usecase"
	"coin-feed/pkg/logger"
	pkgRedis "coin-feed/pkg/redis"
	"coin-feed/pkg/tracing"
	cmcProvider "coin-feed/providers/coinmarketcap"
	redisRepo "coin-feed/repositories/redis"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
)

func Start() {
	logger.InitLogger()
	defer logger.Logger.Sync()

	tp := tracing.InitTracer()
	defer tp.Shutdown(context.Background())

	r := gin.New()
	r.Use(otelgin.Middleware("coin-feed-api"))

	coinMarketCapProvider := cmcProvider.NewProvider(config.UrlCoinMarketCap, config.ApiKeyCoinMarketCap)
	redisClient := pkgRedis.NewRedisClient()
	redisRepository := redisRepo.NewRedisRepository(redisClient)
	fetchCryptoCurrencyMapUseCase := usecase.NewFetchCryptocurrencyMap(coinMarketCapProvider, redisRepository)
	cryptoHandler := NewCryptoHandler(fetchCryptoCurrencyMapUseCase)
	cryptoHandler.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logger.Logger.Info("Starting API server with Gin", zap.String("port", port))
	if err := r.Run(":" + port); err != nil {
		logger.Logger.Fatal("Server failed to start", zap.Error(err))
	}
	logger.Logger.Info("API running with Gin", zap.String("port", port))
}
