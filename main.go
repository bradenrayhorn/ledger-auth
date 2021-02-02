package main

import (
	"github.com/bradenrayhorn/ledger-auth/config"
	"go.uber.org/zap"
)

func main() {
	config.LoadConfig()

	logger := initLogger()
	defer logger.Sync()

	zap.S().Debug("starting ledger-auth service...")
}

func initLogger() *zap.Logger {
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(logger)

	return logger
}
