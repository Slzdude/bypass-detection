package main

import (
	"github.com/Skactor/bypass-detection/config"
	"github.com/Skactor/bypass-detection/logger"
	"github.com/Skactor/bypass-detection/server"
)

func main() {
	err := logger.InitLogger()
	if err != nil {
		logger.Logger.Fatalf("Failed to init logger: %s", err.Error())
	}
	cfg, err := config.Parse("./config.yaml")
	if err != nil {
		logger.Logger.Fatalf("Failed to parse configuration file: %s", err.Error())
		return
	}
	server.StartServer(&cfg.Server)
}
