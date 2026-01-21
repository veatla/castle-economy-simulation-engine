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
	var seed int64 = 1000
	for range 1 {
		world.Agents = append(world.Agents, agents.CreateSimpleAgent(seed, world.Width, world.Height))
	}
	tick := 0
	time := 50 * time.Millisecond
	for range ticker.C {
		tick++
		world.Tick(time)
		// broadcast a lightweight snapshot to connected websocket clients
		BroadcastWorld(tick, &world)
	}
}
