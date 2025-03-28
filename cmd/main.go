package main

import (
	"context"
	"database/sql"
	"google.golang.org/grpc"
	"log"
	"net"
	"simple-service/generated/auth"
	"simple-service/internal/config"
	"simple-service/internal/repo"
	service "simple-service/internal/server"

	"github.com/pkg/errors"
	customLogger "simple-service/internal/logger"
)

func main() {
	var cfg config.AppConfig

	db, err := sql.Open("postgres", "DB_HOST=localhost\nDB_PORT=5433\nDB_NAME=simple_service\nDB_USER=admin\nDB_PASSWORD=admin\nDB_SSL_MODE=disable\nDB_POOL_MAX_CONNS=10\nDB_POOL_MAX_CONN_LIFETIME=300s\nDB_POOL_MAX_CONN_IDLE_TIME=150s")
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	logger, err := customLogger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error initializing logger"))
	}

	userRepo := repo.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, logger)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(errors.Wrap(err, "error initializing listener"))
	}

	grpcServer := grpc.Server{}()
	auth.RegisterAuthServiceServer(grpcServer, authService)

	log.Println("Starting gRPC server on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(errors.Wrap(err, "error starting gRPC server"))
	}
}
