package main

import (
	"./commsclients"
	"./core"
	"./exchangeclients"
	"./handlers"
	log "github.com/cihub/seelog"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	core.LoadEnvVars()
	commsclients.Load()
	exchangeclients.Load()
	commsclients.Telegram.Send(commsclients.TelegramUser, "Bot Started up!")

	go core.CheckPrices()

	router := mux.NewRouter()
	router.HandleFunc("/health", handlers.HealthEndpointHandler).Methods("GET")
	log.Critical(http.ListenAndServe(":5555", router))
}
