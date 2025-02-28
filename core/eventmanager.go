package core

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Event struct represents a process or flow node event
type Event struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Id          string `json:"id"`
	ElementName string `json:"element_name"`
}

type TaskEvent struct {
	Name              string `json:"name"`
	Type              string `json:"type"`
	Id                string `json:"id"`
	ElementName       string `json:"element_name"`
	ProcessInstanceId string `json:"process_instance_id"`
}

// EventManager manages WebSocket clients
type EventManager struct {
	clients  map[*websocket.Conn]bool // Connected clients
	mu       sync.Mutex
	upgrader websocket.Upgrader
}

// NewEventManager creates a new WebSocket manager
func NewEventManager() *EventManager {
	return &EventManager{
		clients: make(map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

// HandleConnections manages new WebSocket connections
func (em *EventManager) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := em.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		return
	}

	em.mu.Lock()
	em.clients[conn] = true
	em.mu.Unlock()

	log.Println("New WebSocket client connected")

	go em.handleMessages(conn)
}

// HandleMessages listens for client disconnects
func (em *EventManager) handleMessages(conn *websocket.Conn) {
	defer func() {
		em.removeClient(conn)
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Client disconnected:", err)
			break
		}
	}
}

// RemoveClient removes a client from the connection pool
func (em *EventManager) removeClient(conn *websocket.Conn) {
	em.mu.Lock()
	defer em.mu.Unlock()
	delete(em.clients, conn)
	log.Println("Client disconnected and removed")
}

// Broadcast sends an event to all connected clients
func (em *EventManager) Broadcast(event interface{}) {
	em.mu.Lock()
	defer em.mu.Unlock()

	message, err := json.Marshal(event)
	if err != nil {
		log.Println("Error encoding event:", err)
		return
	}

	for conn := range em.clients {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("WebSocket send error:", err)
			conn.Close()
			delete(em.clients, conn)
		}
	}
}
