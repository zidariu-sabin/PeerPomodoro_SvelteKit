package main

import (
	"PeerPomodoroBackend_Go/internal/adapters/in_memory"
	"PeerPomodoroBackend_Go/internal/service"
	transportHttp "PeerPomodoroBackend_Go/internal/transport/http"
	wsTransport "PeerPomodoroBackend_Go/internal/transport/websocket"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
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

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	// Initialize repository
	sessionRepo := in_memory.NewSessionRepository()

	// Initialize services
	sessionService := service.NewSessionService(sessionRepo)

	// Initialize HTTP Handler
	httpHandler := transportHttp.NewHandler(sessionService)

	hub := wsTransport.NewHub(sessionService)
	go hub.Run()

	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		wsTransport.ServeWs(hub, w, r, upgrader)
	})

	http.HandleFunc("/create-session", httpHandler.CreateSession)
	http.HandleFunc("/session/", httpHandler.GetSession)

	log.Printf("Server is successfully running on address %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))

}
