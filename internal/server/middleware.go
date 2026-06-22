package server

import (
	"log"
	"net/http"
)

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow your frontend (for now allow all)
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Allowed headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Allowed methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}
