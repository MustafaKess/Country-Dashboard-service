package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// HandleError writes an error response with a specific status code and message.
// It also logs the error.
func HandleError(w http.ResponseWriter, status int, err error, msg string) {
	http.Error(w, msg+": "+err.Error(), status)
	log.Println("Error:", msg, err)
}

// DecodeRequest decodes the JSON request body into the supplied type.
func DecodeRequest[T any](r *http.Request) (T, error) {
	var t T
	err := json.NewDecoder(r.Body).Decode(&t)
	return t, err
}

// Encode encodes data to JSON and writes it to the response with a given status code.
func Encode(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
