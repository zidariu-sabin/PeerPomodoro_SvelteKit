package main

import (
	"PeerPomodoroBackend_Go/internal"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http port number")
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
		if allowedOriginsEnv == "" {
			// For convenience, fallback to a default development origin if the env var is not set.
			log.Println("WARN: ALLOWED_ORIGINS environment variable not set. Defaulting to 'http://localhost:5173'.")
			allowedOriginsEnv = "http://localhost:5173"
		}

		allowedOrigins := strings.Split(allowedOriginsEnv, ",")
		origin := r.Header.Get("Origin")

		for _, allowedOrigin := range allowedOrigins {
			if origin == strings.TrimSpace(allowedOrigin) {
				return true
			}
		}

		log.Printf("Connection from origin '%s' is not allowed.", origin)
		return false
	},
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	hub := internal.NewHub()
	go hub.Run()

	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		internal.ServeWs(hub, w, r, upgrader)
	})

	log.Printf("Server is successfully running on address %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))

}