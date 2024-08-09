package main

import (
	"onlineshop/cmd/handler"
	"onlineshop/internal"

	"onlineshop/internal/service"
	"onlineshop/internal/storage"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @title onlineShop App API
// @version 1.0
// @description API Server for Online Shop Application

// @host localhost:8080
// @BasePath /

// @securityDefinitions.bearer BearerAuth
// @in header
// @name Authorization

func main() {

	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := InitConfig(); err != nil {
		logrus.Fatalf("error to read config %s", err.Error())
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
		logrus.Fatalf("error during connecting to db: %s", err.Error())
	}

	repos := storage.NewStore(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(internal.Server)

	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error until running server: %s", err.Error())
	}

}
