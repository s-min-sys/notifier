package inters

import "github.com/s-min-sys/notifier/pkg"

type Sender interface {
	GetID() pkg.SenderID
	SendTextMessage(message pkg.TextMessage) (err error)

	Wait()
}
