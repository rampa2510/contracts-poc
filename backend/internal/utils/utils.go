package utils

import (
	"encoding/json"
	"net/http"
)

func SendResponse(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(v)
}
