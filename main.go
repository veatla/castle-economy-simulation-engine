package main

import (
	"example/hello/src/agents"
	"example/hello/src/world"
	"fmt"
	"strings"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Millisecond * 50)
	var world = world.NewWorld()
	for range 5 {
		world.Agents = append(world.Agents, agents.CreateSimpleAgent(world.Width, world.Height))

	}
	for range ticker.C {
		world.Tick(0.05)
		fmt.Print("\033[H\033[2J") // очистка экрана
		fmt.Print(render(&world))
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
