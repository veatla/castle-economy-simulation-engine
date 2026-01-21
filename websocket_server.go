package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"example/hello/src/world"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type wsHub struct {
	mu    sync.Mutex
	conns map[*websocket.Conn]bool
}

var hub = &wsHub{conns: make(map[*websocket.Conn]bool)}

func (h *wsHub) addConn(c *websocket.Conn) {
	h.mu.Lock()
	h.conns[c] = true
	h.mu.Unlock()
}

func (h *wsHub) removeConn(c *websocket.Conn) {
	h.mu.Lock()
	delete(h.conns, c)
	h.mu.Unlock()
	c.Close()
}

func (h *wsHub) broadcast(v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		log.Println("ws marshal error:", err)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	for c := range h.conns {
		if err := c.WriteMessage(websocket.TextMessage, b); err != nil {
			log.Println("ws write error, removing conn:", err)
			_ = c.Close()
			delete(h.conns, c)
		}
	}
}

// StartWebSocketServer starts an http server with a /ws endpoint
func StartWebSocketServer() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade error:", err)
			return
		}
		hub.addConn(c)
		// keep reading to detect closed connections
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				hub.removeConn(c)
				return
			}
		}
	})

	log.Println("WebSocket server listening on :8080/ws")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

type agentSnapshot struct {
	ID   int     `json:"id"`
	X    float64 `json:"x"`
	Z    float64 `json:"z"`
	Type string  `json:"type"`
}

// BroadcastMessage is the shape sent to clients
type BroadcastMessage struct {
	Tick    int             `json:"tick"`
	Updated []agentSnapshot `json:"updated"`
}

// BroadcastWorld sends a JSON message with tick and updated agents to clients
func BroadcastWorld(tick int, w *world.World) {
	snap := make([]agentSnapshot, 0, len(w.Agents))
	for _, a := range w.Agents {
		snap = append(snap, agentSnapshot{ID: a.ID, X: a.X, Z: a.Z, Type: "agent"})
	}
	msg := BroadcastMessage{Tick: tick, Updated: snap}
	hub.broadcast(msg)
}
