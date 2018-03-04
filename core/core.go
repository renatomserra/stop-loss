package core

import (
	"./commsclients"
	"./types"
	log "github.com/cihub/seelog"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

var StopLossLimit float64

func CheckPrices() {
	state := types.State{}

	for {
		state.RefreshAll()
		log.Infof("[main.priceCheck] Stop loss set to: %v and current is %v", StopLossLimit, state.CurrentPrice)

		if state.CurrentPrice < StopLossLimit {
			log.Infof("[main.priceCheck] Stop hit, sell all the things")
			go commsclients.Telegram.Send(commsclients.TelegramUser, "Stop Loss triggered!!")

			state.CancelOrders()
			go commsclients.Telegram.Send(commsclients.TelegramUser, "All existing orders cancelled!!")

			state.TriggerStopOrder()
			go commsclients.Telegram.Send(commsclients.TelegramUser, "Sold all the things")
		}

		time.Sleep(15 * time.Second)
	}
}

func StopLoss() float64 {
	value, err := strconv.ParseFloat(os.Getenv("STOP_LOSS_AT"), 64)
	if err != nil {
		panic(err)
	}

	return value
}

func LoadEnvVars() {
	err := godotenv.Load()
	if err != nil {
		log.Error(err)
	}

	StopLossLimit = StopLoss()
}
