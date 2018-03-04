package commsclients

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"os"
	"strconv"
	"time"
)

var TelegramUser *tb.User
var Telegram *tb.Bot

func Load() {
	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("TELEGRAM_TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		panic(err)
	}

	b.Handle("/alive", func(m *tb.Message) {
		b.Send(m.Sender, "Yeah man, im good")
	})

	teleId, err := strconv.Atoi(os.Getenv("TELEGRAM_CHAT_ID"))
	if err != nil {

	}
	Telegram = b
	TelegramUser = &tb.User{
		ID: teleId,
	}
}
