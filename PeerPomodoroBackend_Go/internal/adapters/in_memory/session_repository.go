package in_memory

import (
	"PeerPomodoroBackend_Go/internal/domain"
	"PeerPomodoroBackend_Go/internal/ports"
	"errors"
	"sync"
	"time"
)

type SessionRepository struct {
	sessions map[string]*domain.Session
	mu       sync.RWMutex
}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{
		sessions: make(map[string]*domain.Session),
	}
}

func (r *SessionRepository) Save(session *domain.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[session.ID] = session
	return nil
}

func (r *SessionRepository) Get(id string) (*domain.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	session, exists := r.sessions[id]
	if !exists {
		return nil, ports.ErrSessionNotFound
	}
	return session, nil
}

func (r *SessionRepository) Update(session *domain.Session) error {
	if session == nil {
		return errors.New("session cannot be nil")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.sessions[session.ID]; !exists {
		return ports.ErrSessionNotFound
	}
	r.sessions[session.ID] = session
	return nil
}

func (r *SessionRepository) DeleteInactiveSessions(cutoffTime time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, session := range r.sessions {
		if session.LastActivity.Before(cutoffTime) {
			delete(r.sessions, id)
		}
	}
	return nil
}
