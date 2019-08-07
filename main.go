package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/hvs-fasya/micro/internal/configure"
	"github.com/hvs-fasya/micro/internal/logging"
	"github.com/hvs-fasya/micro/internal/server"
)

const (
	configsDir = "./configs/"
)

func main() {
	err := configure.LoadConfigs(configsDir)
	if err != nil {
		fmt.Printf("load configs error: %s\n", err)
		os.Exit(1)
	}
	configure.Cfg.Server.Logger = logging.SetLoggers()
	defer zap.L().Sync()
	defer configure.Cfg.Server.Logger.Sync()
	server.Run(configure.Cfg.Server)
}
