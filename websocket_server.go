package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"sync"

	"veatla/simulator/src/agents"
	"veatla/simulator/src/constructions"

	"github.com/google/uuid"
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

	log.Println("WebSocket server listening on 127.0.0.1::8080/ws")
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		log.Fatal(err)
	}
}

type agentSnapshot struct {
	ID       uuid.UUID `json:"id"`
	X        float64   `json:"x"`
	Z        float64   `json:"z"`
	Rotation float64   `json:"rotation"`
	Type     string    `json:"type"`
}

type obstacleSnapshot struct {
	ID   uuid.UUID `json:"id"`
	MinX float64   `json:"minX"`
	MinZ float64   `json:"minZ"`
	MaxX float64   `json:"maxX"`
	MaxZ float64   `json:"maxZ"`
	Type string    `json:"type"`
}

// BroadcastMessage is the shape sent to clients
type BroadcastMessage struct {
	Tick      int                `json:"tick"`
	Updated   []agentSnapshot    `json:"updated"`
	Obstacles []obstacleSnapshot `json:"obstacles"`
}

// BroadcastWorld sends only updated agents to clients in batches.
func BroadcastWorld(tick int, updated []agents.Agent, obstacles []constructions.Obstacle) {
	const batchSize = 1000
	total := len(updated)
	obsSnap := make([]obstacleSnapshot, 0, len(obstacles))
	for _, o := range obstacles {
		obsSnap = append(obsSnap, obstacleSnapshot{
			ID:   o.ID,
			MinX: o.MinX,
			MinZ: o.MinZ,
			MaxX: o.MaxX,
			MaxZ: o.MaxZ,
			Type: "obstacle",
		})
	}

	if total == 0 {
		// small heartbeat message
		msg := BroadcastMessage{Tick: tick, Updated: []agentSnapshot{}, Obstacles: obsSnap}
		hub.broadcast(msg)
		return
	}

	for start := 0; start < total; start += batchSize {
		end := start + batchSize
		if end > total {
			end = total
		}
		chunk := updated[start:end]
		snap := make([]agentSnapshot, 0, len(chunk))
		for _, a := range chunk {
			snap = append(snap, agentSnapshot{
				ID:       a.ID,
				X:        a.X,
				Z:        a.Z,
				Type:     "agent",
				Rotation: math.Atan2(a.VZ, a.VX) + math.Pi/2,
			})
		}
		msg := BroadcastMessage{Tick: tick, Updated: snap, Obstacles: obsSnap}
		hub.broadcast(msg)
	}
}
