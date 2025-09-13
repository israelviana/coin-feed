package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() {
	logLevelStr := strings.ToLower(os.Getenv("LOG_LEVEL"))
	var level zapcore.Level
	switch logLevelStr {
	case "debug":
		level = zap.DebugLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(level)
	config.OutputPaths = []string{"stdout", "coin-feed.log"}
	config.ErrorOutputPaths = []string{"stderr"}
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	config.InitialFields = map[string]interface{}{
		"service": "coin-feed",
	}

	config.Sampling = &zap.SamplingConfig{
		Initial:    100,
		Thereafter: 100,
	}

	var err error
	Logger, err = config.Build(zap.AddStacktrace(zap.ErrorLevel))
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
}
