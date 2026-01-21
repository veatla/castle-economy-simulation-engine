package agents

import (
	"example/hello/src/position"
	"math"
	"math/rand"
	"slices"
)

type Agent struct {
	X             float64
	Z             float64
	ID            int
	VX, VZ        float64
	Width, Height float64
	changeDirIn   int
	rng           *rand.Rand
}

var lastID = 0

func CreateSimpleAgent(worldWidth int, seed int64, worldHeight int) Agent {
	id := lastID + 1

	r := rand.New(rand.NewSource(seed + int64(id)))
	pos := position.GetRandomPositionFromWorld(worldWidth, worldHeight)
	angle := r.Float32() * 2 * math.Pi

	return Agent{
		ID:          id,
		X:           float64(pos.X),
		Z:           float64(pos.Z),
		Width:       1.0,
		Height:      1.0,
		VX:          float64(math.Cos(float64(angle))),
		VZ:          float64(math.Sin(float64(angle))),
		changeDirIn: r.Intn(200) + 50,
		rng:         r,
	}
}

func (agent *Agent) Tick() {
	// cell := position.SolidCell{
	// 	X: int(agent.X),
	// 	Z: int(agent.Z),
	// }
	agent.X += agent.VX
	agent.Z += agent.VZ

	agent.changeDirIn--

	if agent.changeDirIn <= 0 {
		angle := agent.rng.Float32() * 2 * math.Pi
		agent.VX = float64(math.Cos(float64(angle)))
		agent.VZ = float64(math.Sin(float64(angle)))
		agent.changeDirIn = agent.rng.Intn(200) + 50
	}

}

func (agent *Agent) AgentWandering(worldWidth, worldHeight int) {

}

func (agent *Agent) MoveAgent(grid map[position.SolidCell][]int, dt float64, worldWidth, worldHeight int) {
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
