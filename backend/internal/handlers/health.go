package handlers

import (
	"encoding/json"
	"net/http"
)

const (
	healthStatusOK = "ok"
)

type healthResponse struct {
	Status string `json:"status"`
}

func Health(w http.ResponseWriter, r *http.Request) {
	response := healthResponse{
		Status: healthStatusOK,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
