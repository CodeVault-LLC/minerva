package logger

import (
	"os"

	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger() (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()

	if os.Getenv("ENV") == "production" {
		cfg = zap.NewProductionConfig()
	}

	cfg.Level.SetLevel(zap.DebugLevel)

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	Log = logger

	return logger, nil
}
