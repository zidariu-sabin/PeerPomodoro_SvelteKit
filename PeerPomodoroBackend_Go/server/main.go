package main

import (
	"PeerPomodoroBackend_Go/internal/adapters/in_memory"
	"PeerPomodoroBackend_Go/internal/middleware"
	"PeerPomodoroBackend_Go/internal/service"
	transportHttp "PeerPomodoroBackend_Go/internal/transport/http"
	wsTransport "PeerPomodoroBackend_Go/internal/transport/websocket"
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var addr = flag.String("addr", "localhost:8080", "http port number")
var upgrader = websocket.Upgrader{
	CheckOrigin: middleware.CheckOrigin,
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

	// Start background cleanup job
	// Run every 10 minutes, delete sessions inactive for 24 hours
	go sessionService.CleanupJob(context.Background(), 10*time.Minute, 24*time.Hour)

	// Initialize HTTP Handler
	httpHandler := transportHttp.NewHandler(sessionService)

	hub := wsTransport.NewHub()
	handler := wsTransport.NewHandler(sessionService, hub)
	go hub.Run()

	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		wsTransport.ServeWs(hub, handler, w, r, upgrader)
	})

	http.HandleFunc("/create-session", middleware.CORS(httpHandler.CreateSession))
	http.HandleFunc("/session/", middleware.CORS(httpHandler.GetSession))

	log.Printf("Server is successfully running on address %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))

}
