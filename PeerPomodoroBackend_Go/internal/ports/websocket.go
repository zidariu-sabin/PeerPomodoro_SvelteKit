package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Connection struct {
	ID     string `json:"id"`
	UserId string `json:"userId"`
}

func initializeConnection(id string, userId string) *Connection {
	return &Connection{
		ID:     id,
		UserId: userId,
	}
}

type ConnectionManager struct {
	Upgrader websocket.Upgrader
}

func NewConnectionManager(upgrader websocket.Upgrader) *ConnectionManager {
	return &ConnectionManager{
		Upgrader: upgrader,
	}
}

func (cm *ConnectionManager) InitializeConnection(w http.ResponseWriter, r *http.Request) {
	// upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := cm.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	log.Printf("Connection from origin '%s' has been upgraded.", r.Host)

	defer log.Printf("Connection from origin '%s' has been closed.", r.Host)
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
