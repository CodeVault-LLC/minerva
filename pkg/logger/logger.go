package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger() (*zap.Logger, error) {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755)
		if err != nil {
			return nil, err
		}
	}

	logFile := "logs/" + time.Now().Format("2006-01-02") + ".log"
	errorLogFile := "logs/" + time.Now().Format("2006-01-02") + "-error.log"

	cfg := zap.Config{
		Development:      true,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout", logFile},
		ErrorOutputPaths: []string{"stderr", errorLogFile, logFile},
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
	}

	if os.Getenv("ENV") == "production" {
		cfg.Development = false
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	Log = logger

	return logger, nil
}
