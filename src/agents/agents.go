package agents

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
	navgrid "veatla/simulator/src/nav-grid"
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
	// Stuck detection
	stuckCounter   int
	lastX, lastZ   float64
	stuckThreshold int
	lastWanderTime time.Time
	lastReplanTick int // cooldown: only replan once per N ticks when stuck
	// Wandering logger
	wanderingEvents []WanderingEvent
	// Path computed by A* for current wandering target
	path      []navgrid.PathPoint
	pathIndex int
}

type Wandering struct {
	X     float64
	Z     float64
	wait  time.Duration
	speed float64
}

// WanderingEvent logs when an agent wanders
type WanderingEvent struct {
	Timestamp        time.Time
	X, Z             float64
	TargetX, TargetZ float64
	Duration         time.Duration
}

// StuckEvent logs when an agent gets stuck
type StuckEvent struct {
	Timestamp time.Time
	X, Z      float64
	Duration  time.Duration
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
		ID:              id,
		X:               tx,
		Z:               tz,
		Width:           1.0,
		Height:          1.0,
		VX:              math.Cos(angle),
		VZ:              math.Sin(angle),
		baseSpeed:       r.Float64()*0.02 + 0.01,
		changeDirIn:     r.Intn(200) + 50,
		rng:             r,
		stuckThreshold:  100,
		lastX:           tx,
		lastZ:           tz,
		lastWanderTime:  time.Now(),
		wanderingEvents: make([]WanderingEvent, 0),
	}
	agent.Wandering = agent.SetWanderingTarget(q)
	// agent.path = nil
	// agent.pathIndex = 0
	return agent
}

func (agent *Agent) Tick(dt time.Duration, q worldQuery.WorldQuery) bool {
	oldX, oldZ := agent.X, agent.Z

	if agent.Wandering.wait <= 0 {
		agent.logWanderingEvent(q)
		agent.Wandering = agent.SetWanderingTarget(q)
	}

	dx := agent.Wandering.X - agent.X
	dz := agent.Wandering.Z - agent.Z
	dist2 := dx*dx + dz*dz
	const reachDist = 0.5

	if dist2 < reachDist*reachDist {
		// Reached target: decrease wait time; only then we may get a new target when wait <= 0
		agent.Wandering.wait -= dt
		return math.Abs(oldX-agent.X) > 1e-9 || math.Abs(oldZ-agent.Z) > 1e-9
	}

	// Move towards wandering target (or stand still if no path)
	agent.MoveTorwardsWanderingTarget(q)
	agent.detectStuck(q)

	// Do NOT decrease wait here: new target only when actually reached or when path is blocked and A* fails

	return math.Abs(oldX-agent.X) > 1e-9 || math.Abs(oldZ-agent.Z) > 1e-9
}

func (agent *Agent) SetWanderingTarget(q worldQuery.WorldQuery) Wandering {
	return agent.setWanderingTargetWithRadius(q, 30.0)
}

// setWanderingTargetWithRadius picks a random target within maxRadius; use smaller radius when stuck so A* can find a path.
func (agent *Agent) setWanderingTargetWithRadius(q worldQuery.WorldQuery, maxRadius float64) Wandering {
	const obstacleOffset = 1.0 // 10-unit offset from buildings/obstacles

	angle := agent.rng.Float64() * 2 * math.Pi
	radius := math.Sqrt(agent.rng.Float64()) * maxRadius

	dx := math.Cos(angle) * radius
	dz := math.Sin(angle) * radius

	tx := agent.X + dx
	tz := agent.Z + dz

	worldWidth, worldHeight := q.GetBoundaries()

	for tx < 0 || tx > worldWidth || tz < 0 || tz > worldHeight || q.IsPointBlocked(tx, tz) {
		angle = agent.rng.Float64() * 2 * math.Pi
		radius = math.Sqrt(agent.rng.Float64()) * maxRadius
		dx = math.Cos(angle) * radius
		dz = math.Sin(angle) * radius
		tx = agent.X + dx
		tz = agent.Z + dz
	}

	// Try to apply obstacle offset
	var targetX, targetZ float64
	if offsetX, offsetZ, ok := utils.FindOffsetPosition(tx, tz, obstacleOffset, q); ok {
		targetX = offsetX
		targetZ = offsetZ
	} else {
		targetX = utils.Clamp(tx, 0, worldWidth)
		targetZ = utils.Clamp(tz, 0, worldHeight)
	}

	// Compute A* path to the chosen target and store on agent
	if path, found, _ := navgrid.AStarPath(agent.X, agent.Z, targetX, targetZ, q, obstacleOffset); found && len(path) > 0 {
		agent.path = path
		// path[0] is start; move toward path[1] first if exists
		if len(path) > 1 {
			agent.pathIndex = 1
		} else {
			agent.pathIndex = 0
		}
	} else {
		// A* failed: don't use a fallback 2-point direct path (it would go through obstacles)
		// Instead, leave path empty so MoveTorwardsWanderingTarget will do direct movement
		// which is safer and will hit the "blocked" case and trigger A* replanning
		agent.path = nil
		agent.pathIndex = 0
	}

	return Wandering{
		speed: 0.03 + agent.rng.Float64()*0.02,
		X:     targetX,
		Z:     targetZ,
		wait:  500 * time.Millisecond,
		// wait: time.Duration(500+agent.rng.Intn(1000)) * time.Millisecond,
	}
}
func (agent *Agent) MoveTorwardsWanderingTarget(q worldQuery.WorldQuery) {
	// If we have a precomputed path, follow it
	if agent.path != nil && agent.pathIndex < len(agent.path) {
		target := agent.path[agent.pathIndex]
		dx := target.X - agent.X
		dz := target.Z - agent.Z
		dist := math.Sqrt(dx*dx + dz*dz)
		// If close enough to waypoint, advance
		if dist < 0.2 {
			agent.pathIndex++
			return
		}

		// Move towards waypoint
		dx /= dist
		dz /= dist
		agent.VX = dx * (agent.baseSpeed + agent.Wandering.speed)
		agent.VZ = dz * (agent.baseSpeed + agent.Wandering.speed)
		nextX := agent.X + agent.VX
		nextZ := agent.Z + agent.VZ

		if !q.IsPointBlocked(nextX, nextZ) {
			worldWidth, worldHeight := q.GetBoundaries()
			agent.X = utils.Clamp(nextX, 0, worldWidth)
			agent.Z = utils.Clamp(nextZ, 0, worldHeight)
			return
		}

		// If the waypoint itself is free, snap to it
		if !q.IsPointBlocked(target.X, target.Z) {
			worldWidth, worldHeight := q.GetBoundaries()
			agent.X = utils.Clamp(target.X, 0, worldWidth)
			agent.Z = utils.Clamp(target.Z, 0, worldHeight)
			agent.pathIndex++
			return
		}

		// Next step blocked: attempt to replan immediately
		agent.navigateWithAStar(q)
		return
	}

	// No path: try moving directly toward target (may help escape tight spots)
	dx := agent.Wandering.X - agent.X
	dz := agent.Wandering.Z - agent.Z
	dist := math.Sqrt(dx*dx + dz*dz)
	if dist < 1e-6 {
		agent.VX = 0
		agent.VZ = 0
		return
	}
	dx /= dist
	dz /= dist
	step := agent.baseSpeed + agent.Wandering.speed
	nextX := agent.X + dx*step
	nextZ := agent.Z + dz*step
	if !q.IsPointBlocked(nextX, nextZ) {
		worldWidth, worldHeight := q.GetBoundaries()
		agent.X = utils.Clamp(nextX, 0, worldWidth)
		agent.Z = utils.Clamp(nextZ, 0, worldHeight)
		agent.VX = dx * step
		agent.VZ = dz * step
		return
	}
	// Direct move blocked: try to find path or new target
	agent.navigateWithAStar(q)
}

// navigateWithAStar uses A* algorithm to find path around obstacles
func (agent *Agent) navigateWithAStar(q worldQuery.WorldQuery) {
	const obstacleOffset = 1.0

	path, found, _ := navgrid.AStarPath(
		agent.X, agent.Z,
		agent.Wandering.X, agent.Wandering.Z,
		q, obstacleOffset,
	)

	if !found || len(path) == 0 {
		// Path blocked / no route: pick a new wandering target instead of staying stuck
		agent.Wandering = agent.SetWanderingTarget(q)
		agent.stuckCounter = 0
		agent.lastX = agent.X
		agent.lastZ = agent.Z
		return
	}

	// Store path for MoveTorwardsWanderingTarget to follow
	agent.path = path
	if len(path) > 1 {
		agent.pathIndex = 1
	} else {
		agent.pathIndex = 0
	}
}

// detectStuck checks if agent is stuck in one location
func (agent *Agent) detectStuck(q worldQuery.WorldQuery) {
	dx := agent.X - agent.lastX
	dz := agent.Z - agent.lastZ
	distance := math.Sqrt(dx*dx + dz*dz)

	const epsilon = 0.01
	const replanCooldown = 50 // only attempt replan once per 50 ticks

	if distance < epsilon {
		agent.stuckCounter++
		if agent.stuckCounter > agent.stuckThreshold {
			log.Printf("Agent %s is STUCK at position (%.2f, %.2f) for %d ticks with distance %.2f",
				agent.ID.String()[:8], agent.X, agent.Z, agent.stuckCounter,
				distance,
			)

			// Attempt to recompute A* path only if cooldown has elapsed (expensive operation)
			if agent.stuckCounter%replanCooldown == 0 {
				const obstacleOffset = 1.0
				if agent.Wandering.X != agent.X || agent.Wandering.Z != agent.Z {
					path, found, _ := navgrid.AStarPath(agent.X, agent.Z, agent.Wandering.X, agent.Wandering.Z, q, obstacleOffset)
					if found && len(path) > 0 {
						agent.path = path
						if len(path) > 1 {
							agent.pathIndex = 1
						} else {
							agent.pathIndex = 0
						}
					} else {
						// Cannot reach current target: try a closer wandering target so A* can find a path
						agent.Wandering = agent.setWanderingTargetWithRadius(q, 8.0)
					}
				}
			}
			agent.stuckCounter = 0
			agent.lastX = agent.X
			agent.lastZ = agent.Z
		}
	} else {
		agent.stuckCounter = 0
		agent.lastX = agent.X
		agent.lastZ = agent.Z
	}
}

// logWanderingEvent logs when agent starts a new wandering target
func (agent *Agent) logWanderingEvent(q worldQuery.WorldQuery) {
	now := time.Now()
	duration := now.Sub(agent.lastWanderTime)

	event := WanderingEvent{
		Timestamp: now,
		X:         agent.X,
		Z:         agent.Z,
		TargetX:   agent.Wandering.X,
		TargetZ:   agent.Wandering.Z,
		Duration:  duration,
	}

	agent.wanderingEvents = append(agent.wanderingEvents, event)
	agent.lastWanderTime = now

	// log.Printf("Agent %s wandering from (%.2f, %.2f) to (%.2f, %.2f) for %v",
	//
	//	agent.ID.String()[:8], event.X, event.Z, event.TargetX, event.TargetZ, event.Duration)
}

// GetWanderingEvents returns all recorded wandering events
func (agent *Agent) GetWanderingEvents() []WanderingEvent {
	return agent.wanderingEvents
}

// ClearWanderingEvents clears the event log
func (agent *Agent) ClearWanderingEvents() {
	agent.wanderingEvents = make([]WanderingEvent, 0)
}

// PrintWanderingLog prints a summary of wandering behavior
func (agent *Agent) PrintWanderingLog() {
	fmt.Printf("\n=== Wandering Log for Agent %s ===\n", agent.ID.String()[:8])
	fmt.Printf("Total wandering events: %d\n", len(agent.wanderingEvents))
	for i, event := range agent.wanderingEvents {
		fmt.Printf("  Event %d: [%s] From (%.2f, %.2f) to (%.2f, %.2f) - Duration: %v\n",
			i+1, event.Timestamp.Format("15:04:05"), event.X, event.Z, event.TargetX, event.TargetZ, event.Duration)
	}
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

// GetPath returns the current computed path
func (agent *Agent) GetPath() []navgrid.PathPoint {
	return agent.path
}

// 	if agent.Z < 0 || agent.Z > float64(worldHeight) {
// 		agent.VZ = -agent.VZ
// 	}
// }
