package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hrygo/council/internal/core/workflow"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all CORS for dev
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Hub maintains the set of active clients and broadcasts messages to the
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan workflow.StreamEvent
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan workflow.StreamEvent),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

// Broadcast sends an event to all connected clients
func (h *Hub) Broadcast(event workflow.StreamEvent) {
	h.broadcast <- event
}

// Client is a middleman between the websocket connection and the hub
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan workflow.StreamEvent
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		message, ok := <-c.send
		if !ok {
			// Hub closed the channel
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}

		if err := json.NewEncoder(w).Encode(message); err != nil {
			return
		}

		if err := w.Close(); err != nil {
			return
		}
	}
}

// HandleWS handles websocket requests from the peer.
func ServeWs(hub *Hub, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan workflow.StreamEvent, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	// Read pump logic can be added if we accept input from WS
}
