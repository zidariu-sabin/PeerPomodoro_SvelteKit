package domain

import "time"

type Session struct {
	ID           string    `json:"id"`
	Clients      []Client  `json:"clients"`
	Timer        *Timer    `json:"timer"`
	LastActivity time.Time `json:"last_activity"`
}


