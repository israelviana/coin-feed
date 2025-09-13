package job

import (
	"context"

	"coin-feed/pkg/logger"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type iSaveLatestCryptoCurrencyUC interface {
	Run(ctx context.Context) error
}

func Start(ctx context.Context, saveLatestCryptoCurrencyUC iSaveLatestCryptoCurrencyUC) {
	c := cron.New()

	_, err := c.AddFunc("0 */6 * * *", func() {
		if err := saveLatestCryptoCurrencyUC.Run(ctx); err != nil {
			logger.Logger.Error("Failed to save latest crypto currency", zap.Error(err))
		} else {
			logger.Logger.Info("Succeed to save latest crypto currency")
		}
	})
	if err != nil {
		logger.Logger.Fatal("Failed to create cron job", zap.Error(err))
	}

	c.Start()

	go func() {
		<-ctx.Done()
		c.Stop()
	}()
}
