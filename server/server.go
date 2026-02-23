package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// StartWebSocketServer starts an HTTP server with a /ws endpoint.
func StartWebSocketServer() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade error:", err)
			return
		}
		hub.addConn(c)
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				hub.removeConn(c)
				return
			}
		}
	})

	log.Println("WebSocket server listening on 127.0.0.1:8080/ws")
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		log.Fatal(err)
	}
}
