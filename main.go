package main

import (
	"github.com/bradenrayhorn/ledger-auth/config"
	"github.com/bradenrayhorn/ledger-auth/database"
	"github.com/bradenrayhorn/ledger-auth/routing"
	"go.uber.org/zap"
)

func main() {
	config.LoadConfig()

	logger := initLogger()
	defer logger.Sync()

	zap.S().Debug("connecting to database...")
	database.Setup()

	zap.S().Debug("connecting to redis...")
	database.SetupRedis()

	zap.S().Debug("starting ledger-auth service...")

	// start http
	r := routing.MakeRouter()

	err := r.Run()

	if err != nil {
		zap.S().Panic(err)
	}
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
