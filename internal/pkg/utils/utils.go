package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ErrorResponse struct {
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Code      int    `json:"code"`
}

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return errors.New("отсутствует тело запроса")
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, payload any) error {
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, ErrorResponse{Message: err.Error(), Code: 0})
}
