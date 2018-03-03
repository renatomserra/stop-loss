package apiclient

import (
	"fmt"
	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/joho/godotenv"
	"os"
)

var KrakenClient = newClient()

func newClient() *krakenapi.KrakenApi {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	return krakenapi.New(os.Getenv("KRAKEN_KEY"), os.Getenv("KRAKEN_SECRET"))
}
