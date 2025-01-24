package main

import (
	"fmt"

	"github.com/Melikhov-p/sso-grpc-go/internal/config"
)

func main() {
	var cfg *config.Config
	cfg = config.MustLoad()

	fmt.Printf("%+v", cfg)

	// TODO: logger

	// TODO: grpc-server

	// TODO: graceful shutdown
}
