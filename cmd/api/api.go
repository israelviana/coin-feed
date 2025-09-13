package api

import (
	"context"
	"os"

	"coin-feed/pkg/logger"
	"coin-feed/pkg/tracing"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
)

func Start(cryptoHandler *CryptoHandler) {
	logger.InitLogger()
	defer logger.Logger.Sync()

	tp := tracing.InitTracer()
	defer tp.Shutdown(context.Background())

	r := gin.New()
	r.Use(otelgin.Middleware("coin-feed-api"))

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
