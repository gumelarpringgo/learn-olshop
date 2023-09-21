package handler

import (
	"encoding/json"
	"net/http"
)

type messageError struct {
	Message string `json:"message"`
}

func WriteErrorResponse(w http.ResponseWriter, code int, err error) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	err = json.NewEncoder(w).Encode(messageError{Message: err.Error()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func WriteDataResponse(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
