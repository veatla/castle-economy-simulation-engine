package world

import (
	"math/rand"

	"veatla/simulator/src/agents"
	"veatla/simulator/src/constructions"
	spatialhash "veatla/simulator/src/spatial-hash"
)

type WorldID string

type World struct {
	rng       *rand.Rand
	Seed      int64
	Width     float64
	Height    float64
	Agents    []agents.Agent
	Obstacles []constructions.Obstacle
	Grid      spatialhash.SpatialHash
}
