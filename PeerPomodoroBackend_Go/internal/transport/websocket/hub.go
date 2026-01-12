package websocket

import (
	"PeerPomodoroBackend_Go/internal/domain"
	"PeerPomodoroBackend_Go/internal/service"
	"encoding/json"
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

					// Broadcast UserLeft to remaining clients
					resp := domain.UserLeftResponse{ClientID: client.id}
					payload, _ := json.Marshal(resp)
					wrapper := domain.Message{Type: domain.MessageTypeUserLeft, Payload: payload}
					finalMsg, _ := json.Marshal(wrapper)

					for remainingClient := range clients {
						select {
						case remainingClient.send <- finalMsg:
						default:
							close(remainingClient.send)
							delete(clients, remainingClient)
						}
					}

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

func (h *Hub) StartTimer(sessionID string) {
	onTick := func(t *domain.Timer) {
		h.broadcastTimerState(sessionID)
	}

	if err := h.sessionService.StartTimer(sessionID, onTick); err != nil {
		log.Printf("Error starting timer for session %s: %v", sessionID, err)
		return
	}
	h.broadcastTimerState(sessionID)
}

func (h *Hub) PauseTimer(sessionID string) {
	if err := h.sessionService.PauseTimer(sessionID); err != nil {
		log.Printf("Error pausing timer for session %s: %v", sessionID, err)
		return
	}
	h.broadcastTimerState(sessionID)
}

func (h *Hub) ResetTimer(sessionID string) {
	if err := h.sessionService.ResetTimer(sessionID); err != nil {
		log.Printf("Error resetting timer for session %s: %v", sessionID, err)
		return
	}
	h.broadcastTimerState(sessionID)
}

func (h *Hub) broadcastTimerState(sessionID string) {
	session, err := h.sessionService.GetSession(sessionID)
	if err != nil {
		log.Printf("Error getting session %s for broadcast: %v", sessionID, err)
		return
	}

	resp := domain.TimerUpdateResponse{Timer: *session.Timer}
	payload, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling timer update: %v", err)
		return
	}

	wrapper := domain.Message{
		Type:    domain.MessageTypeTimerUpdate,
		Payload: payload,
	}
	finalMsg, err := json.Marshal(wrapper)
	if err != nil {
		log.Printf("Error marshalling wrapper: %v", err)
		return
	}

	h.broadcast <- SessionMessage{
		SessionID: sessionID,
		Payload:   finalMsg,
	}
}
