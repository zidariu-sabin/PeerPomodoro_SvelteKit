package service

import (
	"PeerPomodoroBackend_Go/internal/domain"
	"PeerPomodoroBackend_Go/internal/ports"
	"errors"

	"github.com/google/uuid"
)

// SessionManager defines the contract for managing sessions
type SessionManager interface {
	CreateSession(timer domain.Timer) *domain.Session
	GetSession(id string) (*domain.Session, error)
	AddClientToSession(sessionID string, client domain.Client) error
	RemoveClientFromSession(sessionID string, clientID string) error
	StartTimer(sessionID string, onTick func(*domain.Timer)) error
	PauseTimer(sessionID string) error
	ResetTimer(sessionID string) error
}

// SessionService implements SessionManager with different adapters such as in-memory storage/ database
type SessionService struct {
	repo ports.SessionRepository
}

func NewSessionService(repo ports.SessionRepository) *SessionService {
	return &SessionService{
		repo: repo,
	}
}

func (s *SessionService) CreateSession(timer domain.Timer) *domain.Session {
	// Generate a unique UUIDv4
	id := uuid.NewString()

	session := &domain.Session{
		ID:      id,
		Clients: []domain.Client{},
		Timer:   &timer,
	}
	s.repo.Save(session)
	return session
}

func (s *SessionService) GetSession(id string) (*domain.Session, error) {
	return s.repo.Get(id)
}

func (s *SessionService) AddClientToSession(sessionID string, client domain.Client) error {
	session, err := s.repo.Get(sessionID)
	if err != nil {
		return err
	}

	session.Clients = append(session.Clients, client)
	return s.repo.Update(session)
}

func (s *SessionService) RemoveClientFromSession(sessionID string, clientID string) error {
	session, err := s.repo.Get(sessionID)
	if err != nil {
		return err
	}

	for i, client := range session.Clients {
		if client.ID == clientID {
			session.Clients = append(session.Clients[:i], session.Clients[i+1:]...)
			return s.repo.Update(session)
		}
	}
	return errors.New("client not found in session")
}

func (s *SessionService) StartTimer(sessionID string, onTick func(*domain.Timer)) error {
	session, err := s.repo.Get(sessionID)
	if err != nil {
		return err
	}

	session.Timer.OnTick = onTick
	session.Timer.Start()

	return s.repo.Update(session)
}

func (s *SessionService) PauseTimer(sessionID string) error {
	session, err := s.repo.Get(sessionID)
	if err != nil {
		return err
	}

	session.Timer.Pause()
	return s.repo.Update(session)
}

func (s *SessionService) ResetTimer(sessionID string) error {
	session, err := s.repo.Get(sessionID)
	if err != nil {
		return err
	}

	session.Timer.InitializeTimer()
	return s.repo.Update(session)
}
