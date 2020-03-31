package request

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ErrorResponds struct {
	Error string `json:"error"`
}

// SendJSON marshals v to a json struct and sends appropriate headers to w
func SendJSON(w http.ResponseWriter, r *http.Request, v interface{}, statusCode int) error {
	w.Header().Add("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(v)
	if err == nil {
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(v)
		return nil
	}
	log.Print(fmt.Sprintf("Error while encoding JSON: %v", err))
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(ErrorResponds{
		Error: http.StatusText(http.StatusInternalServerError),
	})
	return err
}

func SendServerError(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(ErrorResponds{
		Error: http.StatusText(http.StatusInternalServerError),
	})
}

func SendBadRequestError(w http.ResponseWriter, r *http.Request, message string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(ErrorResponds{
		Error: message,
	})
}
