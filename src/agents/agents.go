package agents

import (
	"example/hello/src/position"

	"github.com/google/uuid"
)

type Agent struct {
	position.Position
	ID     uuid.UUID
	VX, VZ float64
}

func CreateSimpleAgent(worldWidth int, worldHeight int) Agent {

	return Agent{
		ID: uuid.New(),
		Position: position.GetRandomPositionFromWorld(
			worldWidth, worldHeight, 1.0, 1.0,
		),
		VX: 1,
		VZ: 1,
	}
}
