package core

import (
	"github.com/RenatoSerra22/stop-loss/commsclients"
	"github.com/RenatoSerra22/stop-loss/types"
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
		log.Infof("[main.priceCheck] Stop loss set to: %v and current is %v, has position: %v", StopLossLimit, state.CurrentPrice, state.HasPosition())

		if state.CurrentPrice < StopLossLimit && state.HasPosition() {
			log.Infof("[main.priceCheck] Stop hit, sell all the things")
			go commsclients.Telegram.Send(commsclients.TelegramUser, "Stop Loss price reached, cancelling orders!!")

			state.CancelOrders()
			go commsclients.Telegram.Send(commsclients.TelegramUser, "All existing orders cancelled!!, checking for fake dip")

			if shouldSell() {
				state.TriggerStopOrder()
			}
		}

		time.Sleep(15 * time.Second)
	}
}

// Check if it was a quick dip before selling, wait 5 seconds re-check price still below
func shouldSell() bool {
	if os.Getenv("CHECK_FOR_DIP") == "true" {
		time.Sleep(5 * time.Second)

		currentPrice, err := types.FetchPrice()
		if err != nil {
			log.Errorf("[state.isDip] %v", err)
			return false
		} else if currentPrice > StopLossLimit {
			return false
		}
	}

	return true
}

func LoadDefaultStopLoss() float64 {
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

	StopLossLimit = LoadDefaultStopLoss()
}
