package exchangeclients

import (
	krakenapi "github.com/beldur/kraken-go-api-client"
	"os"
)

var Kraken *krakenapi.KrakenApi

func Load() {
	Kraken = krakenapi.New(os.Getenv("KRAKEN_KEY"), os.Getenv("KRAKEN_SECRET"))
}
