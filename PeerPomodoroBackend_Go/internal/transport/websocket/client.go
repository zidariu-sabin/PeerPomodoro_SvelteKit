package websocket

import (
	"PeerPomodoroBackend_Go/internal/domain"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// Client's unique identifier
	id string

	// The session this client belongs to
	sessionID string
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var msg domain.Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Invalid JSON: %v", err)
			continue
		}

		switch msg.Type {
		case domain.MessageTypeJoinSession:
			var req domain.JoinSessionRequest
			if err := json.Unmarshal(msg.Payload, &req); err != nil {
				log.Printf("Invalid JoinSessionRequest: %v", err)
				continue
			}

			// Generate Client ID
			clientID := uuid.NewString()
			dClient := domain.NewClient(clientID, req.UserName)

			// Add to session via service
			if err := c.hub.sessionService.AddClientToSession(req.SessionID, *dClient); err != nil {
				log.Printf("Failed to join session: %v", err)
				
				resp := domain.ErrorResponse{Message: "Session not found or unavailable"}
				payload, _ := json.Marshal(resp)
				msg := domain.Message{Type: domain.MessageTypeError, Payload: payload}
				finalMsg, _ := json.Marshal(msg)
				c.send <- finalMsg
				continue
			}

			c.id = clientID
			c.sessionID = req.SessionID
			c.hub.join <- c

			// Broadcast UserJoined to session
			userJoinedResp := domain.UserJoinedResponse{Client: *dClient}
			userJoinedPayload, _ := json.Marshal(userJoinedResp)
			userJoinedMsg := domain.Message{
				Type:    domain.MessageTypeUserJoined,
				Payload: userJoinedPayload,
			}
			userJoinedFinal, _ := json.Marshal(userJoinedMsg)

			c.hub.broadcast <- SessionMessage{
				SessionID: req.SessionID,
				Payload:   userJoinedFinal,
			}

			// Get updated session data
			session, err := c.hub.sessionService.GetSession(req.SessionID)
			if err != nil {
				log.Printf("Failed to get session: %v", err)
				continue
			}

			// Send confirmation to this client
			response := domain.SessionJoinedResponse{
				SessionID: session.ID,
				ClientID:  clientID,
				Timer:     *session.Timer,
				Clients:   session.Clients,
			}

			respPayload, _ := json.Marshal(response)
			respMsg := domain.Message{
				Type:    domain.MessageTypeSessionJoined,
				Payload: respPayload,
			}

			finalMsg, _ := json.Marshal(respMsg)
			c.send <- finalMsg

		case domain.MessageTypeUpdateUser:
			var req domain.UpdateUserRequest
			if err := json.Unmarshal(msg.Payload, &req); err != nil {
				log.Printf("Invalid UpdateUserRequest: %v", err)
				continue
			}

			if c.sessionID == "" {
				continue
			}

			if err := c.hub.sessionService.UpdateClientName(c.sessionID, c.id, req.Name); err != nil {
				log.Printf("Failed to update user name: %v", err)
				continue
			}

			resp := domain.UserUpdatedResponse{
				ClientID: c.id,
				Name:     req.Name,
			}
			payload, _ := json.Marshal(resp)

			broadcastMsg := domain.Message{
				Type:    domain.MessageTypeUserUpdated,
				Payload: payload,
			}
			finalMsg, _ := json.Marshal(broadcastMsg)

			c.hub.broadcast <- SessionMessage{
				SessionID: c.sessionID,
				Payload:   finalMsg,
			}

		case domain.MessageTypeStartTimer:
			if c.sessionID != "" {
				c.hub.StartTimer(c.sessionID)
			}

		case domain.MessageTypePauseTimer:
			if c.sessionID != "" {
				c.hub.PauseTimer(c.sessionID)
			}

		case domain.MessageTypeStopTimer:
			if c.sessionID != "" {
				c.hub.ResetTimer(c.sessionID)
			}

		default:
			if c.sessionID != "" {
				c.hub.broadcast <- SessionMessage{
					SessionID: c.sessionID,
					Payload:   message,
				}
			}
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, upgrader websocket.Upgrader) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
