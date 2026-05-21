package main

import (
	"github.com/GizmoVault/gotools/base/logx"
	"github.com/s-min-sys/notifier/v2/internal/config"
	"github.com/s-min-sys/notifier/v2/internal/server"
)

func main() {
	logger := logx.NewConsoleLoggerWrapper()
	logger.GetLogger().SetLevel(logx.LevelDebug)

	logger.Info("process start")

	server.NewServer(config.GetConfig(), logger).Wait()
}
