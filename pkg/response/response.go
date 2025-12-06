package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"results,omitempty"`
	Stack   string      `json:"stack,omitempty"`
}

// JSON sends a JSON response with a specific status code
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Success sends a success response
func Success(w http.ResponseWriter, status int, data interface{}) {
	JSON(w, status, data)
}

// Error sends an error response
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, Response{
		Code:    status,
		Message: message,
	})
}