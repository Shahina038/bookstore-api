package handler

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func successResponse(w http.ResponseWriter, status int, message string, data any) {
	writeJSON(w, status, map[string]any{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func errorResponse(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{
		"success": false,
		"message": message,
	})
}
