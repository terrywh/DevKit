package util

import (
	"encoding/json"
	"net/http"
)

type jsonError struct {
	Error string `json:"error"`
}

func JSONError(w http.ResponseWriter, desc string, status int) {
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.Encode(jsonError{Error: desc})
}