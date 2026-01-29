package api

import (
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(HealthResponse{
		Status:  200,
		Message: "Server online",
	})
}
