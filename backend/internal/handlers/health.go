package handlers

import (
	"encoding/json"
	"net/http"
)

const (
	HealthStatusOK = "ok"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

func Health(version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := HealthResponse{
			Status:  HealthStatusOK,
			Version: version,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
