package main

import (
	"time"
	"veatla/simulator/src/agents"
	"veatla/simulator/src/constructions"
	"veatla/simulator/src/world"
)

func main() {
	ticker := time.NewTicker(time.Millisecond * 50)
	w := world.NewWorld(int64(123456), 50, 50)

	go StartWebSocketServer()

	// for range 1 {
	w.Obstacles = append(w.Obstacles,
		constructions.CreateObstacle(1, 1, 10, 10),
		constructions.CreateObstacle(15, 15, 25, 25),
		constructions.CreateObstacle(30, 5, 40, 15),
		constructions.CreateObstacle(40, 40, 50, 50),
	)
	// }

	for i := range w.Obstacles {
		obstacle := &w.Obstacles[i]

		w.Grid.Insert(
			obstacle.ID,
			obstacle.MinX,
			obstacle.MinZ,
			obstacle.MaxX,
			obstacle.MaxZ,
			true,
		)
	}

	for range 1 {
		w.Agents = append(w.Agents, agents.CreateSimpleAgent(&w))
	}

	tick := 0
	tickDur := 50 * time.Millisecond
	for range ticker.C {
		tick++
		w.Grid.Clear(false)

		updated := w.AgentsTick(tickDur)

		// broadcast a lightweight snapshot (only changed agents) to connected websocket clients
		BroadcastWorld(tick, updated, w.Obstacles)
	}
}
