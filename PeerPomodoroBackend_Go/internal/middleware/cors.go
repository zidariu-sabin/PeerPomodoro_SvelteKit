package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"
)

// CheckOrigin verifies if the request origin is allowed based on the ALLOWED_ORIGINS env var.
func CheckOrigin(r *http.Request) bool {
	allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
	if allowedOriginsEnv == "" {
		allowedOriginsEnv = "http://localhost:5173"
	}

	allowedOrigins := strings.Split(allowedOriginsEnv, ",")
	origin := r.Header.Get("Origin")

	// If no origin header is present (e.g. server-to-server or same-origin), strictly speaking 
	// standard CORS might allow it or not depending on policy. 
	// For this app, we'll iterate allowed list.
	
	for _, allowedOrigin := range allowedOrigins {
		if origin == strings.TrimSpace(allowedOrigin) {
			return true
		}
	}
	
	log.Printf("Origin '%s' not allowed", origin)
	return false
}

// CORS wraps an http.HandlerFunc to provide standard CORS headers and checks.
func CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if CheckOrigin(r) {
			origin := r.Header.Get("Origin")
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
