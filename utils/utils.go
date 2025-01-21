package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)


func ParseJSON(r *http.Request, payload interface{}) error {
	if r.Body == nil {
		return errors.New("request body is required")
	}

	// decodes the incoming json into a struct and returns an error if the json is invalid
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		return err
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

var Validate = validator.New()
