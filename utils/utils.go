package utils

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

// GetTokenFromRequest extracts the JWT token from the request.
// It checks the Authorization header for a Bearer token and falls back to the query parameter "token".
func GetTokenFromRequest(r *http.Request) string {
	// Check the Authorization header for a Bearer token
	tokenAuth := r.Header.Get("Authorization")
	if tokenAuth != "" {
		// Extract the token from the "Bearer <token>" format
		parts := strings.Split(tokenAuth, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
	}

	// Fall back to the query parameter "token"
	tokenQuery := r.URL.Query().Get("token")
	if tokenQuery != "" {
		return tokenQuery
	}

	// No token found
	return ""
}
