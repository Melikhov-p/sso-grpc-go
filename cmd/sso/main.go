package main

import (
	"fmt"

	"github.com/Melikhov-p/sso-grpc-go/internal/config"
	"github.com/Melikhov-p/sso-grpc-go/internal/logger"
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

	// TODO: grpc-server

	// TODO: graceful shutdown
}
