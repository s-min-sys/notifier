package main

import (
	"github.com/s-min-sys/notifier/internal/config"
	"github.com/s-min-sys/notifier/internal/server"
	"github.com/sgostarter/i/l"
	"github.com/sgostarter/liblogrus"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := l.NewWrapper(liblogrus.NewLogrusEx(logrus.New()))
	logger.GetLogger().SetLevel(l.LevelDebug)

	logger.Info("process start")

	server.NewServer(config.GetConfig(), logger).Wait()
}
