package main

import (
	"context"

	"github.com/RomanIkonnikov93/niisva/internal/config"
	"github.com/RomanIkonnikov93/niisva/internal/grpcapi"
	"github.com/RomanIkonnikov93/niisva/internal/mkdir"
	"github.com/RomanIkonnikov93/niisva/internal/repository"
	"github.com/RomanIkonnikov93/niisva/internal/server"
	"github.com/RomanIkonnikov93/niisva/pkg/logging"
)

func main() {

	logger := logging.GetLogger()

	cfg, err := config.GetConfig()
	if err != nil {
		logger.Fatalf("GetConfig: %s", err)
	}

	rep, err := repository.NewFileRepository(cfg)
	if err != nil {
		logger.Fatalf("NewFileRepository: %s", err)
	}

	service, err := grpcapi.InitServices(context.Background(), cfg, logger, rep)
	if err != nil {
		logger.Fatalf("InitServices: %s", err)
	}

	err = mkdir.CreateStorageDir(cfg)
	if err != nil {
		logger.Fatalf("СreateStorageDir: %s", err)
	}

	go func() {
		err = service.Run()
		if err != nil {
			logger.Fatalf("Run: %s", err)
		}
	}()

	err = server.StartServer(service, cfg, logger)
	if err != nil {
		logger.Fatalf("StartServer: %s", err)
	}
}
