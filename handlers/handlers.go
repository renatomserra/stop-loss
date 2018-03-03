package handlers

import (
	"net/http"
)

func HealthEndpointHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All good!"))
}
