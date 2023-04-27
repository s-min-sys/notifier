package telebot

import (
	"context"
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (notifier *telebotSenderImpl) userRedisKey(id int64) string {
	return fmt.Sprintf("tele:user:%d", id)
}

func (notifier *telebotSenderImpl) chatRedisKey(id int64) string {
	return fmt.Sprintf("tele:chat:%d", id)
}

func (notifier *telebotSenderImpl) cacheUser(user *tgbotapi.User) (err error) {
	if user == nil {
		return
	}

	key := notifier.userRedisKey(user.ID)

	if _, ok := notifier.dCache.Get(key); ok {
		return
	}

	d, err := json.Marshal(user)
	if err != nil {
		return
	}

	err = notifier.redisCli.Set(context.Background(), key, string(d), 0).Err()

	if err != nil {
		return
	}

	notifier.dCache.SetDefault(key, user.UserName)

	return
}

func (notifier *telebotSenderImpl) cacheChat(chat *tgbotapi.Chat) (err error) {
	if chat == nil {
		return
	}

	key := notifier.chatRedisKey(chat.ID)

	if _, ok := notifier.dCache.Get(key); ok {
		return
	}

	d, err := json.Marshal(chat)
	if err != nil {
		return
	}

	err = notifier.redisCli.Set(context.Background(), key, string(d), 0).Err()

	if err != nil {
		return
	}

	notifier.dCache.SetDefault(key, chat.UserName)

	return
}
