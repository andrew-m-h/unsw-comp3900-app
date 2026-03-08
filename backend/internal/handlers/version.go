package handlers

import (
	"encoding/json"
	"net/http"
)

var AppVersion = "1.0.0"

type versionResponse struct {
	Version string `json:"version"`
}

func Version(w http.ResponseWriter, r *http.Request) {
	response := versionResponse{
		Version: AppVersion,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
