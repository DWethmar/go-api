package common

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// ErrorResponds error responds type
type ErrorResponds struct {
	Error string `json:"error"`
}

// SendJSON set json responds.
func SendJSON(w http.ResponseWriter, r *http.Request, v interface{}, code int) {
	w.Header().Add("Content-Type", "application/json")
	b, err := json.Marshal(v)
	if err != nil {
		log.Print(fmt.Sprintf("Error while encoding JSON: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Internal server error"}`)
	} else {
		w.WriteHeader(code)
		io.WriteString(w, string(b))
	}
}

// SendServerError set server error
func SendServerError(w http.ResponseWriter, r *http.Request) {
	SendJSON(w, r, ErrorResponds{
		Error: http.StatusText(http.StatusInternalServerError),
	}, http.StatusInternalServerError)
}

// SendBadRequestError set bad request responds
func SendBadRequestError(w http.ResponseWriter, r *http.Request, v interface{}) {
	SendJSON(w, r, v, http.StatusBadRequest)
}

// SendNotFoundError set not found responds
func SendNotFoundError(w http.ResponseWriter, r *http.Request) {
	SendJSON(w, r, ErrorResponds{
		Error: "Resource not found.",
	}, http.StatusNotFound)
}
