package main

import (
	"example/hello/src/agents"
	"example/hello/src/world"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Millisecond * 50)
	var world = world.NewWorld(5, 5)
	go StartWebSocketServer()

	for range 2 {
		world.Agents = append(world.Agents, agents.CreateSimpleAgent(world.Width, world.Height))
	}

	tick := 0
	for range ticker.C {
		tick++
		world.Tick(0.05)
		// broadcast a lightweight snapshot to connected websocket clients
		BroadcastWorld(tick, &world)
	}
}
