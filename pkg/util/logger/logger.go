package logger

import (
	"log"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogger() {
	var err error
	Logger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
}

func CloseLogger() {
	if err := Logger.Sync(); err != nil {
		Logger.Fatal("failed to sync logger: %v", zap.String("error", err.Error()))
	}
}
