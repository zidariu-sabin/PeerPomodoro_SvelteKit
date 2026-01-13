package http

import (
	"PeerPomodoroBackend_Go/internal/domain"
	"PeerPomodoroBackend_Go/internal/service"
	"encoding/json"
	"net/http"
	"strings"
)

type Handler struct {
	SessionService service.SessionManager
}

func NewHandler(sessionService service.SessionManager) *Handler {
	return &Handler{
		SessionService: sessionService,
	}
}

func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req domain.CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	timer := domain.NewTimer(req.WorkTime, req.BreakTime, req.TotalRounds)
	session, err := h.SessionService.CreateSession(timer)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}
	resp := domain.SessionCreatedResponse{
		SessionID: session.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) GetSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract Session ID from URL: /session/{id}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}
	sessionID := parts[2]

	session, err := h.SessionService.GetSession(sessionID)
	if err != nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(session); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
