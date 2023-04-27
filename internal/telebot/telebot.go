package telebot

import (
	"github.com/s-min-sys/notifier/pkg"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/patrickmn/go-cache"
	"github.com/s-min-sys/notifier/internal/helper"
	"github.com/s-min-sys/notifier/internal/inters"
	"github.com/sgostarter/i/l"
)

func NewSender(token, apiEndPoint string, redisCli *redis.Client, logger l.Wrapper) inters.Sender {
	if logger != nil {
		logger = l.NewNopLoggerWrapper()
	}

	bot, err := tgbotapi.NewBotAPIWithClient(token, apiEndPoint, helper.NewTelegramBotHTTPClient())
	if err != nil {
		logger.Fatal(err)
	}

	impl := &telebotSenderImpl{
		token:       token,
		apiEndPoint: apiEndPoint,
		redisCli:    redisCli,
		dCache:      cache.New(time.Minute, time.Minute),
		bot:         bot,
	}

	impl.init()

	return impl
}

type telebotSenderImpl struct {
	wg sync.WaitGroup

	logger      l.Wrapper
	token       string
	apiEndPoint string
	redisCli    *redis.Client

	dCache *cache.Cache
	bot    *tgbotapi.BotAPI
}

func (notifier *telebotSenderImpl) GetID() pkg.SenderID {
	return pkg.SenderIDTelegram
}

func (notifier *telebotSenderImpl) Wait() {
	notifier.wg.Wait()
}

func (notifier *telebotSenderImpl) init() {
	notifier.wg.Add(1)

	go notifier.teleRoutine()
}

func (notifier *telebotSenderImpl) teleRoutine() {
	defer notifier.wg.Done()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := notifier.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if update.Message.From != nil {
				_ = notifier.cacheUser(update.Message.From)
			}

			if update.Message.Chat != nil {
				_ = notifier.cacheChat(update.Message.Chat)
			}
		}
	}
}

func (notifier *telebotSenderImpl) SendTextMessage(message pkg.TextMessage) (err error) {
	var chatID int64

	if message.ReceiverType == pkg.ReceiverTypeUser || message.ReceiverType == pkg.ReceiverTypeGroup {
		chatID, err = strconv.ParseInt(message.Receiver, 10, 64)
		if err != nil {
			return
		}
	}

	_, err = notifier.bot.Send(tgbotapi.NewMessage(chatID, message.Text))

	return
}
