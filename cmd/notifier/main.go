package main

import (
	"github.com/s-min-sys/notifier/internal/server"
	"github.com/s-min-sys/notifier/pkg"
	"github.com/sgostarter/i/l"
	"github.com/sgostarter/liblogrus"
	"github.com/sirupsen/logrus"
	"time"
)

func main() {
	logger := l.NewWrapper(liblogrus.NewLogrusEx(logrus.New()))
	logger.GetLogger().SetLevel(l.LevelDebug)

	logger.Info("process start")

	s := server.NewServer(logger)
	time.Sleep(time.Second * 10)
	_ = s.SendTextMessage(pkg.TextMessage{
		SenderID:     pkg.SenderIDTelegram,
		ReceiverType: pkg.ReceiverTypeGroup,
		Receiver:     "-986630020",
		Text:         "notifier started",
	})

	s.Wait()
}
