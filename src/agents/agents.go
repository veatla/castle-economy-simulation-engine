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
	baseSpeed     float64
	rng           *rand.Rand
	Wandering
}

type Wandering struct {
	X     float64
	Z     float64
	wait  time.Duration
	speed float64
}

var lastID = 0

func uuidToInt64(u uuid.UUID) uint64 {
	return binary.BigEndian.Uint64(u[8:16])
}
func CreateSimpleAgent(seed int64, worldWidth, worldHeight float64) Agent {
	lastID += 1
	id := uuid.New()
	r := rand.New(rand.NewSource(seed + -int64(uuidToInt64(id))))
	angle := r.Float64() * 2 * math.Pi

	agent := Agent{
		ID:          id,
		X:           r.Float64() * worldWidth,
		Z:           r.Float64() * worldHeight,
		Width:       1.0,
		Height:      1.0,
		VX:          math.Cos(angle),
		VZ:          math.Sin(angle),
		baseSpeed:   r.Float64()*0.02 + 0.01,
		changeDirIn: r.Intn(200) + 50,
		rng:         r,
	}
	agent.Wandering = agent.SetWanderingTarget(worldWidth, worldHeight)
	return agent
}
func clamp(value, min, max float64) float64 {
	return math.Min(math.Max(value, min), max)
}
func (agent *Agent) Tick(dt time.Duration, worldWidth, worldHeight float64) bool {
	oldX, oldZ := agent.X, agent.Z

	if agent.Wandering.wait <= 0 {
		agent.Wandering = agent.SetWanderingTarget(worldWidth, worldHeight)
	}

	dx := agent.Wandering.X - agent.X
	dz := agent.Wandering.Z - agent.Z
	dist2 := dx*dx + dz*dz
	const reachDist = 0.5

	if dist2 < reachDist*reachDist {
		agent.Wandering.wait -= dt
		// no movement in this tick
		return math.Abs(oldX-agent.X) > 1e-9 || math.Abs(oldZ-agent.Z) > 1e-9
	}

	agent.MoveTorwardsWanderingTarget()
	return math.Abs(oldX-agent.X) > 1e-9 || math.Abs(oldZ-agent.Z) > 1e-9
}

func (agent *Agent) SetWanderingTarget(worldWidth, worldHeight float64) Wandering {

	const maxRadius = 30.0

	angle := agent.rng.Float64() * 2 * math.Pi

	// sqrt — обязательно, иначе будет bias к центру
	radius := math.Sqrt(agent.rng.Float64()) * maxRadius

	dx := math.Cos(angle) * radius
	dz := math.Sin(angle) * radius

	tx := agent.X + dx
	tz := agent.Z + dz

	if tx < 0 || tx > worldWidth || tz < 0 || tz > worldHeight {
		// если вышли — просто пробуем ещё раз
		return agent.SetWanderingTarget(worldWidth, worldHeight)
	}

	return Wandering{
		speed: 0.03 + agent.rng.Float64()*0.02,
		X:     clamp(agent.X+dx, 0, worldWidth),
		Z:     clamp(agent.Z+dz, 0, worldHeight),
		wait:  time.Duration(500+agent.rng.Intn(1200)) * time.Millisecond,
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

	agent.VX = dx * (agent.baseSpeed + agent.Wandering.speed)
	agent.VZ = dz * (agent.baseSpeed + agent.Wandering.speed)

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
