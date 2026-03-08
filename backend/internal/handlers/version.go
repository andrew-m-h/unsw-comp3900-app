package handlers

import (
	"encoding/json"
	"net/http"
)

type versionResponse struct {
	Version string `json:"version"`
}

func Version(version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := versionResponse{
			Version: version,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
