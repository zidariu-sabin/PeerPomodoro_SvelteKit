package domain

import "encoding/json"

// Message defines the structure for all incoming and outgoing websocket messages.
// The `Type` field is used to determine how to process the `Payload`.
type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}