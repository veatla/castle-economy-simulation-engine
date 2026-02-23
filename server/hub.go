package server

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

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
