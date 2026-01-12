package http

import (
	"PeerPomodoroBackend_Go/internal/domain"
	"PeerPomodoroBackend_Go/internal/service"
	"encoding/json"
	"net/http"
	"os"
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
	// CORS handling
	allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
	if allowedOriginsEnv == "" {
		allowedOriginsEnv = "http://localhost:5173"
	}

	allowedOrigins := strings.Split(allowedOriginsEnv, ",")
	origin := r.Header.Get("Origin")
	allow := false

	for _, o := range allowedOrigins {
		if strings.TrimSpace(o) == origin {
			allow = true
			break
		}
	}

	if allow {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

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
	session := h.SessionService.CreateSession(*timer)

	resp := domain.SessionCreatedResponse{
		SessionID: session.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) GetSession(w http.ResponseWriter, r *http.Request) {
	// CORS handling
	allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
	if allowedOriginsEnv == "" {
		allowedOriginsEnv = "http://localhost:5173"
	}

	allowedOrigins := strings.Split(allowedOriginsEnv, ",")
	origin := r.Header.Get("Origin")
	allow := false

	for _, o := range allowedOrigins {
		if strings.TrimSpace(o) == origin {
			allow = true
			break
		}
	}

	if allow {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

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

	_, err := h.SessionService.GetSession(sessionID)
	if err != nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
