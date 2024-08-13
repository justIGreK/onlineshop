package main

import (
	"onlineshop/cmd/handler"
	"onlineshop/internal"

	"onlineshop/internal/service"
	"onlineshop/internal/storage"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// @title onlineShop App API
// @version 1.0
// @description API Server for Online Shop Application

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() {
		if syncErr := logger.Sync(); syncErr != nil {
			logger.Error("Failed to sync logger", zap.Error(syncErr))
		}
	}()
	if err := InitConfig(); err != nil {
		logger.Fatal("error to read config %s", zap.String("error", err.Error()))
	}
	db, err := storage.NewPostgresDB(storage.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logger.Fatal("error during connecting to db: %s", zap.String("error", err.Error()))
	}
	repos := storage.NewStore(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(internal.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logger.Fatal("error until running server: %s", zap.String("error", err.Error()))
	}
}
