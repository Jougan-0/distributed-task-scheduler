package api

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Client holds one WebSocket connection plus a channel to send messages out
type Client struct {
	conn *websocket.Conn
	send chan string
}

type Hub struct {
	mu      sync.Mutex
	clients map[*Client]bool
}

// global hub instance
var wsHub = &Hub{
	clients: make(map[*Client]bool),
}

var upgrader = websocket.Upgrader{
	// CheckOrigin just accepts everything for local dev
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WsHandler is called when the client hits "/ws"
func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan string, 256),
	}
	wsHub.register(client)

	// Start goroutines to read/write
	go client.writePump()
	go client.readPump()
}

// BroadcastLog can be called to send `message` to all connected clients
func BroadcastLog(message string) {
	wsHub.broadcast(message)
}

// register a new client
func (h *Hub) register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c] = true
	log.Println("New WebSocket client connected:", c.conn.RemoteAddr())
}

// unregister client
func (h *Hub) unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[c]; ok {
		delete(h.clients, c)
		close(c.send)
	}
}

// broadcast message to all connected clients
func (h *Hub) broadcast(message string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for c := range h.clients {
		select {
		case c.send <- message:
			// success
		default:
			// client is not receiving, close it
			close(c.send)
			delete(h.clients, c)
		}
	}
}

// readPump reads inbound messages (if you want to handle messages from clients)
func (c *Client) readPump() {
	defer func() {
		c.close()
	}()
	for {
		// We don't necessarily care about incoming data, so we just read to avoid blocking
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("WS read error: %v", err)
			break
		}
		log.Printf("Got message from client: %s", string(msg))
	}
}

// writePump sends messages from c.send channel to the WebSocket
func (c *Client) writePump() {
	for msg := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Printf("WS write error: %v", err)
			break
		}
	}
	c.close()
}

func (c *Client) close() {
	c.conn.Close()
	wsHub.unregister(c)
}
