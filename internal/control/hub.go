package control

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // All origins allowed — LAN-only tool, no cloud relay.
	},
}

// Hub maintains connected WebSocket clients and routes messages between
// remotes (mobile PWA) and players (TV browser).
type Hub struct {
	mu           sync.RWMutex
	clients      map[*Client]bool
	register     chan *Client
	unregister   chan *Client
	done         chan struct{}
	onConnect    func(Role)
	onDisconnect func(Role)
}

// NewHub creates a new Hub. Call Run() in a goroutine to start processing.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client, 16),
		unregister: make(chan *Client, 16),
		done:       make(chan struct{}),
	}
}

// Run processes register, unregister, and shutdown events.
// It blocks until Stop is called.
func (h *Hub) Run() {
	for {
		select {
		case <-h.done:
			h.mu.Lock()
			for c := range h.clients {
				close(c.send)
				delete(h.clients, c)
			}
			h.mu.Unlock()
			return
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c] = true
			onConnect := h.onConnect
			h.mu.Unlock()
			if onConnect != nil {
				onConnect(c.role)
			}
			log.Printf("ws: %s connected (%s)", c.role, c.conn.RemoteAddr())
		case c := <-h.unregister:
			h.mu.Lock()
			existed := false
			if _, ok := h.clients[c]; ok {
				existed = true
				delete(h.clients, c)
				close(c.send)
			}
			onDisconnect := h.onDisconnect
			h.mu.Unlock()
			if existed && onDisconnect != nil {
				onDisconnect(c.role)
			}
			log.Printf("ws: %s disconnected (%s)", c.role, c.conn.RemoteAddr())
		}
	}
}

// SetLifecycleHooks configures callbacks for client connect/disconnect events.
func (h *Hub) SetLifecycleHooks(onConnect, onDisconnect func(Role)) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.onConnect = onConnect
	h.onDisconnect = onDisconnect
}

// Stop shuts down the Hub, closing all client connections.
func (h *Hub) Stop() {
	close(h.done)
}

// route delivers a message from src to all clients of the opposite role.
// Remote messages go to players; player messages go to remotes.
func (h *Hub) route(src *Client, msg []byte) {
	var target Role
	switch src.role {
	case RoleRemote:
		target = RolePlayer
	case RolePlayer:
		target = RoleRemote
	default:
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients {
		if c.role == target {
			select {
			case c.send <- msg:
			default:
				// Slow consumer — drop message. Ping/pong timeout will clean up.
			}
		}
	}
}

// HandleWS upgrades an HTTP request to a WebSocket connection.
// The client role is set via the "role" query parameter ("remote" or "player").
func (h *Hub) HandleWS(w http.ResponseWriter, r *http.Request) {
	roleParam := r.URL.Query().Get("role")
	role := Role(roleParam)
	if role != RoleRemote && role != RolePlayer {
		http.Error(w,
			fmt.Sprintf("invalid role %q: use %q or %q", roleParam, RoleRemote, RolePlayer),
			http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ws: upgrade error: %v", err)
		return
	}

	client := &Client{
		hub:  h,
		conn: conn,
		send: make(chan []byte, sendBufferSize),
		role: role,
	}

	h.register <- client
	go client.writePump()
	go client.readPump()
}

// ClientCount returns the number of currently connected clients.
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
