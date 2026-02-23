package main

import (
	"time"

	"veatla/simulator/server"
	"veatla/simulator/src/agents"
	"veatla/simulator/src/constructions"
	"veatla/simulator/src/world"
)

func main() {
	ticker := time.NewTicker(50 * time.Millisecond)
	w := world.NewWorld(123456, 50, 50)

	go server.StartWebSocketServer()

	w.Obstacles = append(w.Obstacles,
		constructions.CreateObstacle(1, 1, 10, 10),
		constructions.CreateObstacle(15, 15, 25, 25),
		constructions.CreateObstacle(30, 5, 40, 15),
		constructions.CreateObstacle(40, 40, 50, 50),
	)

	for i := range w.Obstacles {
		obstacle := &w.Obstacles[i]
		w.Grid.Insert(obstacle.ID, obstacle.MinX, obstacle.MinZ, obstacle.MaxX, obstacle.MaxZ, true)
	}

	for range 1 {
		w.Agents = append(w.Agents, agents.CreateSimpleAgent(&w))
	}

	tickDur := 50 * time.Millisecond
	tick := 0
	for range ticker.C {
		tick++
		w.Grid.Clear(false)
		updated := w.AgentsTick(tickDur)
		server.BroadcastWorld(tick, updated, w.Obstacles)
	}
}
