package job

import (
	"context"

	"coin-feed/pkg/logger"

	"go.uber.org/zap"
)

type iSaveLatestCryptoCurrencyUC interface {
	Run(ctx context.Context) error
}

func Start(parentCtx context.Context, uc iSaveLatestCryptoCurrencyUC) {
	err := uc.Run(parentCtx)
	if err != nil {
		logger.Logger.Error("Failed to save latest crypto currency", zap.Error(err))
	}

	logger.Logger.Info("Succeed to save latest crypto currency")
}
