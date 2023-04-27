package pkg

type SenderID string

const (
	SenderIDAll      SenderID = "all"
	SenderIDTelegram SenderID = "telegram"
	SenderIDWeChat   SenderID = "wechat"
)

type ReceiverType int

const (
	ReceiverTypeLogicUser ReceiverType = iota
	ReceiverTypeUser
	ReceiverTypeGroup
)

type TextMessage struct {
	SenderID     SenderID
	ReceiverType ReceiverType
	Receiver     string
	Text         string
}
