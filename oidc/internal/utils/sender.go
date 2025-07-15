package utils

import (
	"encoding/json"
	"net/http"
)

func ResponseSend(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Invalid data", http.StatusInternalServerError)
	}
}
