package websocket

import (
	"PeerPomodoroBackend_Go/internal/service"
	"log"
)

// SessionMessage represents a message to be broadcasted to a specific session
type SessionMessage struct {
	SessionID string
	Payload   []byte
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients in the same session.
type Hub struct {
	// Registered sessions and their clients: SessionID -> Client -> true
	sessions map[string]map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan SessionMessage

	// Register requests from the clients (initial connection).
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Session join requests.
	join chan *Client

	sessionService service.SessionManager
}

func NewHub(ss service.SessionManager) *Hub {
	return &Hub{
		broadcast:      make(chan SessionMessage),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		join:           make(chan *Client),
		sessions:       make(map[string]map[*Client]bool),
		sessionService: ss,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// Client connected but hasn't joined a session yet
			log.Printf("Client connected: %p", client)
		case client := <-h.unregister:
			if client.sessionID != "" {
				if clients, ok := h.sessions[client.sessionID]; ok {
					delete(clients, client)
					if len(clients) == 0 {
						delete(h.sessions, client.sessionID)
					}
					// Remove from domain session
					h.sessionService.RemoveClientFromSession(client.sessionID, client.id)
				}
			}
			close(client.send)
		case client := <-h.join:
			if _, ok := h.sessions[client.sessionID]; !ok {
				h.sessions[client.sessionID] = make(map[*Client]bool)
			}
			h.sessions[client.sessionID][client] = true
			log.Printf("Client %s joined session %s", client.id, client.sessionID)
		case message := <-h.broadcast:
			if clients, ok := h.sessions[message.SessionID]; ok {
				for client := range clients {
					select {
					case client.send <- message.Payload:
					default:
						close(client.send)
						delete(clients, client)
					}
				}
				if len(clients) == 0 {
					delete(h.sessions, message.SessionID)
				}
			}
		}
	}
}
