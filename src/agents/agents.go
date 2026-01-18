package agents

import (
	"example/hello/src/position"
)

type Agent struct {
	position.Position
	ID     int
	VX, VZ float32
}

var lastID = 0

func CreateSimpleAgent(worldWidth float32, worldHeight float32) Agent {
	lastID += 1

	return Agent{
		ID: lastID,
		Position: position.GetRandomPositionFromWorld(
			worldWidth, worldHeight, 1.0, 1.0,
		),
		VX: 1,
		VZ: 1,
	}
}
