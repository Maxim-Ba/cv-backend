package apierrors

import (
	"encoding/json"
	"errors"
	"net/http"
)

// ErrNotFound is a sentinel error indicating a resource was not found.
var ErrNotFound = errors.New("not found")

// APIError is the standard error response body for all API errors.
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// WriteError writes a structured JSON error response.
func WriteError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIError{ //nolint:errcheck
		Code:    status,
		Message: message,
	})
}

// IsNotFound reports whether err wraps ErrNotFound.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}
