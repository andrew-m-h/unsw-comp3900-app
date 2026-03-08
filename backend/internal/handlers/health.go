package handlers

import (
	"encoding/json"
	"net/http"
)

const (
	healthStatusOK = "ok"
)

type healthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

func Health(version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := healthResponse{
			Status:  healthStatusOK,
			Version: version,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
