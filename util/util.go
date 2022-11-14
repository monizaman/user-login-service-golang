package util

import (
	"encoding/json"
	"net/http"
)

// ResponseWithError general response error
func ResponseWithError(w http.ResponseWriter, code int, message string) {
	ResponseWithJSON(w, code, map[string]string{"error": message})
}

// ResponseWithJSON general response
func ResponseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
