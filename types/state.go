package types

import (
	"../apiclient"
	krakenapi "github.com/beldur/kraken-go-api-client"
	log "github.com/cihub/seelog"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

type State struct {
	currentPrice float64
	openOrders   []string
	balance      Balance
}

func (s *State) RefreshAll() {
	wg.Add(3)

	go s.refreshPrice()
	go s.refreshBalance()
	go s.refreshOrders()

	wg.Wait()
}

func (s *State) refreshPrice() {
	defer wg.Done()

	ticker, err := apiclient.KrakenClient.Ticker(krakenapi.XXBTZEUR)
	if err != nil {
		log.Errorf("[state.refresh.ticker] %v", err)
	}
	currentPrice, err := strconv.ParseFloat(ticker.XXBTZEUR.Bid[0], 64)
	if err != nil {
		log.Errorf("[state.refresh.parsePrice] %v", err)
	}
	s.currentPrice = currentPrice
}

func (s *State) refreshBalance() {
	defer wg.Done()

	balance, err := apiclient.KrakenClient.Balance()
	if err != nil {
		log.Errorf("[state.refresh.balance] %v", err)
	}
	s.balance = Balance{
		BTC: round_balance(balance.XXBT),
		EUR: round_balance(balance.ZEUR),
	}
}

func (s *State) refreshOrders() {
	defer wg.Done()

	orders, err := apiclient.KrakenClient.OpenOrders(nil)
	if err != nil {
		log.Errorf("[state.refresh.orders] %v", err)
	}

	var orderList []string
	for k, _ := range orders.Open {
		orderList = append(orderList, k)
	}

	s.openOrders = orderList
}

func round_balance(value float64) float64 {
	if value > 0.0001 {
		return value
	}

	return 0
}
