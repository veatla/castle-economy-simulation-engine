package main

import (
	"example/hello/src/agents"
	"example/hello/src/world"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Millisecond * 50)
	w := world.NewWorld(5, 5)
	worldSeed := int64(123456)

	go StartWebSocketServer()

	for range 1 {
		w.Agents = append(w.Agents, agents.CreateSimpleAgent(worldSeed, w.Width, w.Height))
	}

	tick := 0
	tickDur := 50 * time.Millisecond
	for range ticker.C {
		tick++
		updated := w.Tick(tickDur)
		// broadcast a lightweight snapshot (only changed agents) to connected websocket clients
		BroadcastWorld(tick, updated)
	}
}
