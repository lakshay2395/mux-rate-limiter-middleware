package middleware

import (
	"encoding/json"
	"net/http"
)

func TooManyRequests(w http.ResponseWriter) {
	w.WriteHeader(http.StatusTooManyRequests)
	w.Header().Set("Content-Type", "application/json")
	data := map[string]string{}
	data["message"] = "Requests limit execeeded"
	payload, _ := json.Marshal(data)
	w.Write(payload)
}
