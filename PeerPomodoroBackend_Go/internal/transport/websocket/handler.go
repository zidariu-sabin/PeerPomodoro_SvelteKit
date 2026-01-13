package websocket

import (
	"PeerPomodoroBackend_Go/internal/domain"
	"PeerPomodoroBackend_Go/internal/service"
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

// Handler manages WebSocket message processing and orchestrates business logic
type Handler struct {
	sessionService service.SessionManager
	hub            *Hub
}

func NewHandler(ss service.SessionManager, hub *Hub) *Handler {
	return &Handler{
		sessionService: ss,
		hub:            hub,
	}
}

// HandleMessage processes an incoming message from a client
func (h *Handler) HandleMessage(client *Client, msg domain.Message) {
	switch msg.Type {
	case domain.MessageTypeJoinSession:
		h.handleJoinSession(client, msg.Payload)
	case domain.MessageTypeUpdateUser:
		h.handleUpdateUser(client, msg.Payload)
	case domain.MessageTypeStartTimer:
		h.handleStartTimer(client)
	case domain.MessageTypePauseTimer:
		h.handlePauseTimer(client)
	case domain.MessageTypeStopTimer:
		h.handleResetTimer(client)
	default:
		// Default behavior: broadcast raw message to session
		if client.sessionID != "" {
			h.hub.broadcast <- SessionMessage{
				SessionID: client.sessionID,
				Payload:   mustMarshal(msg), // Re-marshal the original message
			}
		}
	}
}

func (h *Handler) handleJoinSession(client *Client, payload []byte) {
	var req domain.JoinSessionRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		log.Printf("Invalid JoinSessionRequest: %v", err)
		return
	}

	// Generate Client ID
	clientID := uuid.NewString()
	dClient := domain.NewClient(clientID, req.UserName)

	// Add to session via service
	if err := h.sessionService.AddClientToSession(req.SessionID, *dClient); err != nil {
		log.Printf("Failed to join session: %v", err)

		resp := domain.ErrorResponse{Message: "Session not found or unavailable"}
		h.sendToClient(client, domain.MessageTypeError, resp)
		return
	}

	client.id = clientID
	client.sessionID = req.SessionID
	h.hub.join <- client

	// Broadcast UserJoined to session
	h.broadcastToSession(req.SessionID, domain.MessageTypeUserJoined, domain.UserJoinedResponse{Client: *dClient})

	// Get updated session data
	session, err := h.sessionService.GetSession(req.SessionID)
	if err != nil {
		log.Printf("Failed to get session: %v", err)
		return
	}

	// Send confirmation to this client
	response := domain.SessionJoinedResponse{
		SessionID: session.ID,
		ClientID:  clientID,
		Timer:     session.Timer,
		Clients:   session.Clients,
	}
	h.sendToClient(client, domain.MessageTypeSessionJoined, response)
}

func (h *Handler) handleUpdateUser(client *Client, payload []byte) {
	var req domain.UpdateUserRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		log.Printf("Invalid UpdateUserRequest: %v", err)
		return
	}

	if client.sessionID == "" {
		return
	}

	if err := h.sessionService.UpdateClientName(client.sessionID, client.id, req.Name); err != nil {
		log.Printf("Failed to update user name: %v", err)
		return
	}

	h.broadcastToSession(client.sessionID, domain.MessageTypeUserUpdated, domain.UserUpdatedResponse{
		ClientID: client.id,
		Name:     req.Name,
	})
}

func (h *Handler) handleStartTimer(client *Client) {
	if client.sessionID == "" {
		return
	}

	// Define the callback for timer ticks
	onTick := func(t *domain.Timer) {
		h.broadcastTimerState(client.sessionID)
	}

	if err := h.sessionService.StartTimer(client.sessionID, onTick); err != nil {
		log.Printf("Error starting timer: %v", err)
		return
	}
	h.broadcastTimerState(client.sessionID)
}

func (h *Handler) handlePauseTimer(client *Client) {
	if client.sessionID == "" {
		return
	}
	if err := h.sessionService.PauseTimer(client.sessionID); err != nil {
		log.Printf("Error pausing timer: %v", err)
		return
	}
	h.broadcastTimerState(client.sessionID)
}

func (h *Handler) handleResetTimer(client *Client) {
	if client.sessionID == "" {
		return
	}
	if err := h.sessionService.ResetTimer(client.sessionID); err != nil {
		log.Printf("Error resetting timer: %v", err)
		return
	}
	h.broadcastTimerState(client.sessionID)
}

func (h *Handler) HandleDisconnect(client *Client) {
	if client.sessionID != "" && client.id != "" {
		if err := h.sessionService.RemoveClientFromSession(client.sessionID, client.id); err != nil {
			log.Printf("Error removing client from session: %v", err)
		}
	}
}

// Helpers

func (h *Handler) broadcastTimerState(sessionID string) {
	session, err := h.sessionService.GetSession(sessionID)
	if err != nil {
		log.Printf("Error getting session for broadcast: %v", err)
		return
	}

	h.broadcastToSession(sessionID, domain.MessageTypeTimerUpdate, domain.TimerUpdateResponse{Timer: session.Timer})
}

func (h *Handler) sendToClient(client *Client, msgType string, payload interface{}) {
	payloadBytes, _ := json.Marshal(payload)
	msg := domain.Message{Type: msgType, Payload: payloadBytes}
	finalMsg, _ := json.Marshal(msg)
	client.send <- finalMsg
}

func (h *Handler) broadcastToSession(sessionID string, msgType string, payload interface{}) {
	payloadBytes, _ := json.Marshal(payload)
	msg := domain.Message{Type: msgType, Payload: payloadBytes}
	finalMsg, _ := json.Marshal(msg)
	h.hub.broadcast <- SessionMessage{SessionID: sessionID, Payload: finalMsg}
}

func mustMarshal(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
