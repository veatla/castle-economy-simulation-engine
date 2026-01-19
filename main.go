package main

import (
	"example/hello/src/agents"
	"example/hello/src/world"
	"strings"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Millisecond * 50)
	var world = world.NewWorld(5, 5)
	// start websocket server to broadcast simulation state
	go StartWebSocketServer()
	for range 1 {
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
func render(world *world.World) string {
	const W, H = 50, 50
	grid := make([][]rune, H)

	for i := range grid {
		grid[i] = make([]rune, W)
		for j := range grid[i] {
			grid[i][j] = '.'
		}
	}

	for _, a := range world.Agents {
		x := int(a.X / world.Width * W)
		z := int(a.Z / world.Height * H)

		if x >= 0 && x < W && z >= 0 && z < H {
			grid[z][x] = '@'
		}
	}

	var sb strings.Builder
	for _, row := range grid {
		sb.WriteString(string(row))
		sb.WriteRune('\n')
	}
	return sb.String()
}
