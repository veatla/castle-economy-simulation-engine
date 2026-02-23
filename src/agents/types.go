package agents

import (
	"math/rand"
	"time"
	navgrid "veatla/simulator/src/nav-grid"

	"github.com/google/uuid"
)

// pathState holds A* path and follow index
type pathState struct {
	path      []navgrid.PathPoint
	pathIndex int
}

// stuckState holds data for stuck detection and replan cooldown
type stuckState struct {
	counter        int
	lastX, lastZ   float64
	threshold      int
	lastReplanTick int
}

// wanderingLog holds event log and last wander time
type wanderingLog struct {
	events         []WanderingEvent
	lastWanderTime time.Time
}

// Agent is the main agent type: position, velocity, wandering target and internal state.
type Agent struct {
	X, Z         float64
	ID           uuid.UUID
	VX, VZ       float64
	Width, Height float64
	changeDirIn  int
	baseSpeed    float64
	rng          *rand.Rand
	Wandering

	path pathState
	stuck stuckState
	log  wanderingLog

	// NoPath is set when A* fails to find a path to current target (used by websocket etc.)
	NoPath bool
}

// GetPath returns the current computed path.
func (a *Agent) GetPath() []navgrid.PathPoint { return a.path.path }

// Wandering is the current wandering target and timing.
type Wandering struct {
	X     float64
	Z     float64
	wait  time.Duration
	speed float64
}

// WanderingEvent logs when an agent wanders.
type WanderingEvent struct {
	Timestamp        time.Time
	X, Z             float64
	TargetX, TargetZ float64
	Duration         time.Duration
}

// StuckEvent logs when an agent gets stuck.
type StuckEvent struct {
	Timestamp time.Time
	X, Z      float64
	Duration  time.Duration
}

