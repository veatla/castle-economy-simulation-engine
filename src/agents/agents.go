package agents

import (
	"encoding/binary"
	"math"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Agent struct {
	X             float64
	Z             float64
	ID            uuid.UUID
	VX, VZ        float64
	Width, Height float64
	changeDirIn   int
	rng           *rand.Rand
	Wandering
}

type Wandering struct {
	X    float64
	Z    float64
	wait time.Duration
}

var lastID = 0

func uuidToInt64(u uuid.UUID) int64 {
	return int64(binary.BigEndian.Uint64(u[8:16]))
}
func CreateSimpleAgent(seed int64, worldWidth, worldHeight int) Agent {
	lastID += 1
	id := uuid.New()

	r := rand.New(rand.NewSource(seed + uuidToInt64(id)))
	angle := r.Float32() * 2 * math.Pi
	agent := Agent{
		ID:          id,
		X:           r.Float64() * float64(worldWidth),
		Z:           r.Float64() * float64(worldHeight),
		Width:       1.0,
		Height:      1.0,
		VX:          float64(math.Cos(float64(angle))),
		VZ:          float64(math.Sin(float64(angle))),
		changeDirIn: r.Intn(200) + 50,
		rng:         r,
	}
	agent.Wandering = agent.SetWanderingTarget(worldWidth, worldHeight)
	return agent
}
func clamp(value, min, max float64) float64 {
	return math.Min(math.Max(value, min), max)
}
func (agent *Agent) Tick(dt time.Duration, worldWidth, worldHeight int) {
	if int(agent.X) == int(agent.Wandering.X) && int(agent.Z) == int(agent.Wandering.Z) {
		agent.Wandering.wait -= dt
	} else {
		agent.MoveTorwardsWanderingTarget()
	}
	if agent.Wandering.wait <= 0 {
		agent.Wandering = agent.SetWanderingTarget(worldWidth, worldHeight)
	}
}

func (agent *Agent) SetWanderingTarget(worldWidth, worldHeight int) Wandering {
	radius := 20.0

	dx := (agent.rng.Float64()*2 - 1) * radius
	dz := (agent.rng.Float64()*2 - 1) * radius

	return Wandering{
		X:    clamp(agent.X+dx, 0, float64(worldWidth)),
		Z:    clamp(agent.Z+dz, 0, float64(worldHeight)),
		wait: 1000 * time.Millisecond,
	}
}
func (agent *Agent) MoveTorwardsWanderingTarget() {
	dx := agent.Wandering.X - agent.X
	dz := agent.Wandering.Z - agent.Z

	length := math.Sqrt(dx*dx + dz*dz)
	if length == 0 {
		agent.VX = 0
		agent.VZ = 0
		return
	}

	dx /= length
	dz /= length

	agent.VX = dx * 0.05
	agent.VZ = dz * 0.05

	agent.X += agent.VX
	agent.Z += agent.VZ
}

// func (agent *Agent) WanderingTarget(worldWidth, worldHeight int) {
// 	agent.Wandering = Wandering{
// 		X:    agent.rng.Float64() * float64(worldWidth),
// 		Z:    agent.rng.Float64() * float64(worldHeight),
// 		wait: time.Duration.Milliseconds(500),
// 	}
// }

// func (agent *Agent) MoveAgent(grid map[position.SolidCell][]int, dt float64, worldWidth, worldHeight int) {
// 	var cell = position.SolidCell{
// 		X: int(agent.X),
// 		Z: int(agent.Z),
// 	}
// 	agent.X = agent.X + agent.VX*dt
// 	agent.Z = agent.Z + agent.VZ*dt

// 	var newCell = position.SolidCell{
// 		X: int(agent.X),
// 		Z: int(agent.Z),
// 	}

// 	if newCell != cell {
// 		index := slices.Index(grid[cell], agent.ID)
// 		if index != -1 {
// 			grid[cell] = slices.Delete(grid[cell], index, index+1)
// 		}
// 	}

// 	if agent.X < 0 || agent.X > float64(worldWidth) {
// 		agent.VX = -agent.VX
// 	}

// 	if agent.Z < 0 || agent.Z > float64(worldHeight) {
// 		agent.VZ = -agent.VZ
// 	}
// }
