package v1

import (
	"encoding/json"
	"net/http"

	"github.com/zqz/web/backend/internal/handler"
)

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

// JSON writes a JSON response
func JSON(w http.ResponseWriter, status int, data interface{}) {
	handler.SetContentType(w, handler.ContentTypeJSON)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Error writes an error JSON response
func Error(w http.ResponseWriter, status int, err error) {
	JSON(w, status, ErrorResponse{
		Error:   err.Error(),
		Message: http.StatusText(status),
	})
}

// ErrorMessage writes an error JSON response with custom message
func ErrorMessage(w http.ResponseWriter, status int, message string) {
	JSON(w, status, ErrorResponse{
		Error:   message,
		Message: http.StatusText(status),
	})
}
