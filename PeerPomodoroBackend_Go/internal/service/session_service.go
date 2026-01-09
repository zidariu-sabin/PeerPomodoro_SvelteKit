package service

import (
	"PeerPomodoroBackend_Go/internal/domain"
	"errors"
	"sync"
)

// SessionManager defines the contract for managing sessions
type SessionManager interface {
	CreateSession(id string, timer domain.Timer) *domain.Session
	GetSession(id string) (*domain.Session, error)
	AddClientToSession(sessionID string, client domain.Client) error
	RemoveClientFromSession(sessionID string, clientID string) error
}

// SessionService implements SessionManager with in-memory storage
type SessionService struct {
	sessions map[string]*domain.Session
	mu       sync.RWMutex
}

func NewSessionService() *SessionService {
	return &SessionService{
		sessions: make(map[string]*domain.Session),
	}
}

func (s *SessionService) CreateSession(id string, timer domain.Timer) *domain.Session {
	s.mu.Lock()
	defer s.mu.Unlock()

	session := &domain.Session{
		ID:      id,
		Clients: []domain.Client{},
		Timer:   timer,
	}
	s.sessions[id] = session
	return session
}

func (s *SessionService) GetSession(id string) (*domain.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[id]
	if !exists {
		return nil, errors.New("session not found")
	}
	return session, nil
}

func (s *SessionService) AddClientToSession(sessionID string, client domain.Client) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return errors.New("session not found")
	}

	session.Clients = append(session.Clients, client)
	return nil
}

func (s *SessionService) RemoveClientFromSession(sessionID string, clientID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return errors.New("session not found")
	}

	for i, client := range session.Clients {
		if client.ID == clientID {
			session.Clients = append(session.Clients[:i], session.Clients[i+1:]...)
			return nil
		}
	}
	return errors.New("client not found in session")
}
