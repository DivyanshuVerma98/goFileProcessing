package handlers

import (
	"encoding/json"
	"net/http"
)

func SendResponse(w http.ResponseWriter, msg string, status int, data interface{}) {
	var response = map[string]interface{}{
		"status":  status,
		"message": msg,
		"data":    data,
	}
	// w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
