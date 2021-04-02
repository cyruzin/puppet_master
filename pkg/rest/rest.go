package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

// APIError is used for handling the API errors.
type APIError struct {
	Error  string `json:"error"`
	Status int    `json:"status,omitempty"`
}

// APIMessage is used for handling the API success messages.
type APIMessage struct {
	Message string `json:"message"`
	Status  int    `json:"status,omitempty"`
}

// EncodeJSON encodes the payload and returns a JSON response.
func EncodeJSON(w http.ResponseWriter, httpStatus int, payload interface{}) {
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(payload)
}

// DecodeJSON decodes the JSON payload.
func DecodeJSON(r io.Reader, payload interface{}) error {
	if err := json.NewDecoder(r).Decode(payload); err != nil {
		return err
	}

	return nil
}

// ParseID parses the given ID to int64.
func ParseID(sid string) (int64, error) {
	fmt.Println(sid)
	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// EncodeError encodes and logs the error payload.
func EncodeError(w http.ResponseWriter, r *http.Request, err error, code int) {
	log.Error().
		Err(err).
		Stack().
		Int("status", code).
		Str("method", r.Method).
		Str("end-point", r.RequestURI).
		Msg(err.Error())

	w.WriteHeader(code)

	e := &APIError{Error: err.Error(), Status: code}

	if err := json.NewEncoder(w).Encode(e); err != nil {
		return
	}
}

// EncodeErrorGraphql encodes and logs the error payload.
func EncodeErrorGraphql(w http.ResponseWriter, r *http.Request, err error) {
	log.Error().
		Err(err).
		Stack().
		Msg(err.Error())

	e := &APIError{Error: err.Error()}

	if err := json.NewEncoder(w).Encode(e); err != nil {
		return
	}
}
