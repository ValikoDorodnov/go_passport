package rest

import (
	"encoding/json"
	"io"
	"net/http"
)

func ResponseOk(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func ResponseErrors(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(ErrorResponse{Errors: []string{err.Error()}})
}

func ParseRequestBody(body io.ReadCloser, data interface{}) error {
	return json.NewDecoder(body).Decode(&data)
}
