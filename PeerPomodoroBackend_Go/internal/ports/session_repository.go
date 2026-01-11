package ports

import (
	"PeerPomodoroBackend_Go/internal/domain"
	"errors"
)

var (
	ErrSessionNotFound = errors.New("session not found")
)

type SessionRepository interface {
	Save(session *domain.Session) error
	Get(id string) (*domain.Session, error)
	Update(session *domain.Session) error
}
