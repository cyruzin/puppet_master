package util

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

// AuthError is used for handling the API errors.
type AuthError struct {
	Error string `json:"error"`
}

// DecodeError logs the error and returns a JSON response.
func DecodeError(w http.ResponseWriter, r *http.Request, err error) {
	log.Error().Msg(err.Error())

	e := &AuthError{Error: err.Error()}

	if err := json.NewEncoder(w).Encode(e); err != nil {
		w.Write([]byte("could not encode the payload"))
		return
	}
}
