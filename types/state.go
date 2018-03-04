package types

import (
	"./exchangeclients"
	krakenapi "github.com/beldur/kraken-go-api-client"
	log "github.com/cihub/seelog"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

type State struct {
	CurrentPrice float64
	OpenOrders   []string
	Balance      Balance
}

// Public method to refresh all data async
func (s *State) RefreshAll() {
	wg.Add(3)

	go s.refreshPrice()
	go s.refreshBalance()
	go s.refreshOrders()

	wg.Wait()
}

// Get the current ticker price
func (s *State) refreshPrice() {
	defer wg.Done()

	ticker, err := exchangeclients.Kraken.Ticker(krakenapi.XXBTZEUR)
	if err != nil {
		log.Errorf("[state.refresh.ticker] %v", err)
	}
	currentPrice, err := strconv.ParseFloat(ticker.XXBTZEUR.Bid[0], 64)
	if err != nil {
		log.Errorf("[state.refresh.parsePrice] %v", err)
	}
	s.CurrentPrice = currentPrice
}

// Get the latest account balance
func (s *State) refreshBalance() {
	defer wg.Done()

	balance, err := exchangeclients.Kraken.Balance()
	if err != nil {
		log.Errorf("[state.refresh.balance] %v", err)
	}
	if balance == nil {
		return
	}
	s.Balance = Balance{
		BTC: balance.XXBT,
		EUR: balance.ZEUR,
	}
}

// Get the list of current open orders so they can be cancelled
func (s *State) refreshOrders() {
	defer wg.Done()

	orders, err := exchangeclients.Kraken.OpenOrders(nil)
	if err != nil {
		log.Errorf("[state.refresh.orders] %v", err)
	}
	if orders == nil {
		s.OpenOrders = []string{}

		return
	}

	var orderList []string
	for k, _ := range orders.Open {
		orderList = append(orderList, k)
	}

	s.OpenOrders = orderList
}

// Cancel all open orders
func (s *State) CancelOrders() {
	for _, order := range s.OpenOrders {
		_, err := exchangeclients.Kraken.CancelOrder(order)
		if err != nil {
			log.Errorf("[state.cancel.order] %v", err)
		}
	}
}

// Send market sell
func (s *State) TriggerStopOrder() {
	_, err := exchangeclients.Kraken.AddOrder("XXBTZEUR", "sell", "market", strconv.FormatFloat(s.Balance.BTC, 'f', 4, 64), nil)
	if err != nil {
		log.Errorf("[state.triggerstop] %v", err)
	}
}
