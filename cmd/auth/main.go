package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// CORS Middleware
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		headers := w.Header()
		headers.Add("Access-Control-Allow-Credentials", "true")
		headers.Add("Access-Control-Allow-Origin", "http://localhost:3000")
		headers.Add("Vary", "Origin")
		headers.Add("Vary", "Access-Control-Request-Method")
		headers.Add("Vary", "Access-Control-Request-Headers")
		headers.Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token")
		headers.Add("Access-Control-Allow-Methods", "GET, POST,OPTIONS")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		fmt.Println("ok")

		// Next
		next.ServeHTTP(w, r)
		return
	})
}

func main() {
	fmt.Println("Starting API.")

	r := mux.NewRouter()

	r.Use(CORS)

	r.HandleFunc("/signin", Signin)
	r.HandleFunc("/welcome", Welcome)
	r.HandleFunc("/refresh", Refresh)

	http.Handle("/", r)

	// Apply the CORS middleware to our top-level router, with the defaults.
	http.ListenAndServe(":8000", r)
}
