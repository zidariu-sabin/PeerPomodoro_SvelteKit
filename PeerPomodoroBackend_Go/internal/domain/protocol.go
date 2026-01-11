package domain

// Incoming Message Types
const (
	MessageTypeJoinSession = "join_session"
)

// Outgoing Message Types
const (
	MessageTypeSessionCreated = "session_created"
	MessageTypeSessionJoined  = "session_joined"
	MessageTypeUserJoined     = "user_joined"
	MessageTypeUserLeft       = "user_left"
	MessageTypeError          = "error"
)

// HTTP
// CreateSessionRequest is the payload for creating a new session via HTTP
type CreateSessionRequest struct {
	WorkTime    int64 `json:"work_time"`
	BreakTime   int64 `json:"break_time"`
	TotalRounds int64 `json:"total_rounds"`
}

// websocket

// SessionCreatedResponse is the payload sent back when a session is created
type SessionCreatedResponse struct {
	SessionID string `json:"session_id"`
}

// JoinSessionRequest is the payload for joining a session
type JoinSessionRequest struct {
	SessionID string `json:"session_id"`
	UserName  string `json:"user_name"`
}

// SessionJoinedResponse is the payload sent back when a user joins a session
type SessionJoinedResponse struct {
	SessionID string   `json:"session_id"`
	ClientID  string   `json:"client_id"`
	Timer     Timer    `json:"timer"`
	Clients   []Client `json:"clients"`
}

// UserJoinedResponse is broadcast when a new user joins the session
type UserJoinedResponse struct {
	Client Client `json:"client"`
}

// UserLeftResponse is broadcast when a user leaves the session
type UserLeftResponse struct {
	ClientID string `json:"client_id"`
}

// ErrorResponse is a generic error payload
type ErrorResponse struct {
	Message string `json:"message"`
}
