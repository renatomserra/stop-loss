package types

import (
	"fmt"
	"github.com/RenatoSerra22/stop-loss/commsclients"
	"github.com/RenatoSerra22/stop-loss/exchangeclients"
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

	currentPrice, err := FetchPrice()
	if err != nil {
		log.Errorf("[state.refresh.price] %v", err)
		return
	}
	s.CurrentPrice = currentPrice
}

// Get the latest account balance
func (s *State) refreshBalance() {
	defer wg.Done()

	balance, err := exchangeclients.Kraken.Balance()
	if err != nil {
		log.Errorf("[state.refresh.balance] %v", err)
		return
	}

	if s.Balance.BTC != balance.XXBT || s.Balance.EUR != balance.ZEUR { //Balance has changed, update
		message := fmt.Sprintf("New Balance BTC: %v, EUR: %v€", balance.XXBT, balance.ZEUR)
		go commsclients.Telegram.Send(commsclients.TelegramUser, message)
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
		return
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
			return
		}
	}
}

// Cancel all open orders
func (s *State) HasPosition() bool {
	if s.Balance.BTC > 0.00001 {
		return true
	}
	return false
}

// Send market sell
func (s *State) TriggerStopOrder() {
	///Kraken rounds numbers making this fail, it will fail to place the order because we dont own the rounded balance
	sell_volume := strconv.FormatFloat(s.Balance.BTC-0.00002, 'f', 5, 64)
	_, err := exchangeclients.Kraken.AddOrder("XXBTZEUR", "sell", "market", sell_volume, nil)
	if err != nil {
		log.Errorf("[state.triggerstop] %v ", err)
		message := fmt.Sprintf("Failed to sell with volume %v because: %v", sell_volume, err)
		go commsclients.Telegram.Send(commsclients.TelegramUser, message)
		return
	}
	go commsclients.Telegram.Send(commsclients.TelegramUser, "Sold all the things")
}

func FetchPrice() (float64, error) {
	ticker, err := exchangeclients.Kraken.Ticker(krakenapi.XXBTZEUR)
	if err != nil {
		return 0, err
	}
	currentPrice, err := strconv.ParseFloat(ticker.XXBTZEUR.Bid[0], 64)
	if err != nil {
		return 0, err
	}
	return currentPrice, nil
}
