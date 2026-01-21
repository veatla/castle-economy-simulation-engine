package agents

import (
	"example/hello/src/position"
	"slices"

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

func (agent *Agent) AgentWandering(worldWidth int, worldHeight int) {

}

func (agent *Agent) MoveAgent(grid map[position.SolidCell][]uuid.UUID, dt float64, worldWidth, worldHeight int) {
	var cell = position.SolidCell{
		X: int(agent.X),
		Z: int(agent.Z),
	}
	agent.X = agent.X + agent.VX*dt
	agent.Z = agent.Z + agent.VZ*dt

	var newCell = position.SolidCell{
		X: int(agent.X),
		Z: int(agent.Z),
	}

	if newCell != cell {
		index := slices.Index(grid[cell], agent.ID)
		if index != -1 {
			grid[cell] = slices.Delete(grid[cell], index, index+1)
		}
	}

	if agent.X < 0 || agent.X > float64(worldWidth) {
		agent.VX = -agent.VX
	}

	if agent.Z < 0 || agent.Z > float64(worldHeight) {
		agent.VZ = -agent.VZ
	}
}
