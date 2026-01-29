package agents

import (
	"math"
	"math/rand"
	"time"
	"veatla/simulator/src/utils"
	worldQuery "veatla/simulator/src/world-query"

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

func CreateSimpleAgent(q worldQuery.WorldQuery) Agent {
	id := uuid.New()
	r := rand.New(rand.NewSource(q.GetWorldSeed() + utils.UUIDToInt64(id)))
	angle := r.Float64() * 2 * math.Pi
	worldWidth, worldHeight := q.GetBoundaries()

	tx := r.Float64() * worldWidth
	tz := r.Float64() * worldHeight

	for q.IsPointBlocked(tx, tz) {
		tx = r.Float64() * worldWidth
		tz = r.Float64() * worldHeight
	}

	agent := Agent{
		ID:          id,
		X:           tx,
		Z:           tz,
		Width:       1.0,
		Height:      1.0,
		VX:          math.Cos(angle),
		VZ:          math.Sin(angle),
		baseSpeed:   r.Float64()*0.02 + 0.01,
		changeDirIn: r.Intn(200) + 50,
		rng:         r,
	}
	agent.Wandering = agent.SetWanderingTarget(q)
	return agent
}

func (agent *Agent) Tick(dt time.Duration, q worldQuery.WorldQuery) bool {
	oldX, oldZ := agent.X, agent.Z

	if agent.Wandering.wait <= 0 {
		agent.Wandering = agent.SetWanderingTarget(q)
	}

	dx := agent.Wandering.X - agent.X
	dz := agent.Wandering.Z - agent.Z
	dist2 := dx*dx + dz*dz
	const reachDist = 0.5

	if dist2 < reachDist*reachDist {
		agent.Wandering.wait -= dt
		return math.Abs(oldX-agent.X) > 1e-9 || math.Abs(oldZ-agent.Z) > 1e-9
	}

	agent.MoveTorwardsWanderingTarget(q)
	return math.Abs(oldX-agent.X) > 1e-9 || math.Abs(oldZ-agent.Z) > 1e-9
}

func (agent *Agent) SetWanderingTarget(q worldQuery.WorldQuery) Wandering {

	const maxRadius = 30.0

	angle := agent.rng.Float64() * 2 * math.Pi

	radius := math.Sqrt(agent.rng.Float64()) * maxRadius

	dx := math.Cos(angle) * radius
	dz := math.Sin(angle) * radius

	tx := agent.X + dx
	tz := agent.Z + dz

	worldWidth, worldHeight := q.GetBoundaries()

	if tx < 0 || tx > worldWidth || tz < 0 || tz > worldHeight || q.IsPointBlocked(tx, tz) {
		return agent.SetWanderingTarget(q)
	}

	return Wandering{
		speed: 0.03 + agent.rng.Float64()*0.02,
		X:     utils.Clamp(agent.X+dx, 0, worldWidth),
		Z:     utils.Clamp(agent.Z+dz, 0, worldHeight),
		wait:  time.Duration(500+agent.rng.Intn(1200)) * time.Millisecond,
	}
}
func (agent *Agent) MoveTorwardsWanderingTarget(q worldQuery.WorldQuery) {
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
	nextZ := agent.Z + agent.VZ
	nextX := agent.X + agent.VX

	if q.IsPointBlocked(nextX, nextZ) {
		agent.ChooseNewDirection(q.RandomFloat())
	} else {
		agent.X += agent.VX
		agent.Z += agent.VZ
	}
}

func (agent *Agent) ChooseNewDirection(randomFloat float64) {

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
