package main

import (
	"fmt"
	"net"

	"github.com/Melikhov-p/sso-grpc-go/internal/config"
	"github.com/Melikhov-p/sso-grpc-go/internal/grpc"
	"github.com/Melikhov-p/sso-grpc-go/internal/logger"
	"github.com/Melikhov-p/sso-grpc-go/internal/service"
	"github.com/Melikhov-p/sso-grpc-go/internal/storage/file"
	"go.uber.org/zap"
)

func main() {
	var cfg *config.Config
	cfg = config.MustLoad()

	fmt.Printf("%+v\n", cfg)

	log, err := logger.SetupLogger(cfg.Env)
	if err != nil {
		panic("error setup logger " + err.Error())
	}

	log.Debug("logger is running")

	if cfg.Storage.StorageMode != config.FileMode {
		panic("Database storage mode is not working yet.")
	}
	storage := file.NewFileStorage(log, cfg.Storage.FileStorage.FilePath)

	// TODO: grpc-server
	server := grpc.NewServer(log)
	authService := service.NewAuthService(log, storage, storage, storage)

	server.RegisterAuthService(authService)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	if err != nil {
		log.Fatal("net.Listen() error", zap.Error(err))
	}

	if err = server.RPC.Serve(l); err != nil {
		log.Fatal("error serving listener", zap.Error(err))
	}

	log.Debug("server is running", zap.Int("PORT", cfg.GRPC.Port))

	// TODO: graceful shutdown
}
