package main

import (
	"github.com/go-redis/redis"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"onlineshop/cmd/handler"
	"onlineshop/internal"
	"onlineshop/internal/service"
	"onlineshop/internal/storage"
	"onlineshop/internal/ws"
	grpcrequest "onlineshop/pkg/grpcReq"
	"onlineshop/pkg/publisher"
	"onlineshop/pkg/util/logger"
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
	logger.InitLogger()
	defer logger.CloseLogger()

	if err := InitConfig(); err != nil {
		logger.Logger.Fatal("error to read config %s", zap.String("error", err.Error()))
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
		logger.Logger.Fatal("error during connecting to db: %s", zap.String("error", err.Error()))
	}
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		logger.Logger.Info("nats is not connected")
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Logger.Fatal("Failed to connect to gRPC server: %v", zap.String("error", err.Error()))
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logger.Logger.Fatal("failed to close connection: %v", zap.String("error", err.Error()))
		}
	}()
	grpcSender := grpcrequest.NewGrpcRequst(conn)
	repos := storage.NewStore(db)
	natsSender := publisher.NewNATSMessageSender(nc, repos.UserList)
	userList := service.NewUserService(repos.UserList)
	shopHandler := handler.NewHandler(
		service.NewAuthService(repos.Authorization, *grpcSender),
		userList,
		service.NewProdService(repos.Product),
		service.NewCartService(natsSender, repos.Cart, repos.Product, repos.Order, repos.UserList),
		service.NewOrderService(repos.Order),
	)
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	hub := ws.NewHub(rdb)
	wsHandler := ws.NewHandler(hub, userList)
	go hub.Run()
	srv := new(internal.Server)
	if err := srv.Run(viper.GetString("port"), handler.InitRoutes(shopHandler, wsHandler)); err != nil {
		logger.Logger.Fatal("error until running server: %s", zap.String("error", err.Error()))
	}
}
