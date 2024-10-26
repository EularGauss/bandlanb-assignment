package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

// ErrorResponse defines the structure of the error response.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code) // Set response status code
	json.NewEncoder(w).Encode(ErrorResponse{
		Code:    code,
		Message: message,
	})
}

func generateID() string {
	return uuid.New().String()
}
