package main

import (
	"github.com/RenatoSerra22/stop-loss/commsclients"
	"github.com/RenatoSerra22/stop-loss/core"
	"github.com/RenatoSerra22/stop-loss/exchangeclients"
	"github.com/RenatoSerra22/stop-loss/handlers"
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
