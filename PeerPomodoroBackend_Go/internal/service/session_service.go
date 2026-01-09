package service

import (
	"PeerPomodoroBackend_Go/internal/domain"
)

type SessionManager interface {
	CreateSession(id string, timer domain.Timer) *domain.Session
	AddClientToSession(session *domain.Session, client domain.Client)
	RemoveClientFromSession(session *domain.Session, clientID string)
}

// Create a database store that will manage sessions and clients
// use key-value pairs
// for now it will be held in memory
func newSession(id string, timer domain.Timer) *domain.Session {
	return &domain.Session{
		ID:      id,
		Clients: []domain.Client{},
		Timer:   timer,
	}
}

func addClientToSession(session *domain.Session, client domain.Client) {
	session.Clients = append(session.Clients, client)
}

func removeClientFromSession(session *domain.Session, clientID string) {
	for i, client := range session.Clients {
		if client.ID == clientID {
			session.Clients = append(session.Clients[:i], session.Clients[i+1:]...)
			break
		}
	}
}
