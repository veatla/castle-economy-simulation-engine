package agents

import (
	"example/hello/src/position"

	"github.com/google/uuid"
)

type Agent struct {
	X             float64
	Z             float64
	ID            uuid.UUID
	VX, VZ        float64
	Width, Height float64
}

func CreateSimpleAgent(worldWidth int, worldHeight int) Agent {

	pos := position.GetRandomPositionFromWorld(worldWidth, worldHeight)

	return Agent{
		ID:     uuid.New(),
		X:      float64(pos.X),
		Z:      float64(pos.Z),
		Width:  1.0,
		Height: 1.0,
		VX:     1,
		VZ:     1,
	}
}

func AgentWandering(agent *Agent, worldWidth int, worldHeight int) {

}
