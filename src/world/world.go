package world

import (
	"math/rand"
	"sync"
	"time"
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

func (w *World) IsPointBlocked(x, z float64) bool {
	return w.Grid.IsPointBlocked(x, z)
}

func (w *World) RandomFloat() float64 {
	return w.rng.Float64()
}

func (w *World) GetBoundaries() (width, height float64) {
	return w.Width, w.Height
}

func (w *World) ApplyBoundaries(a *agents.Agent) {
	if a.X < 0 {
		a.X = 0
		a.VX = -a.VX
	}
	if a.X > float64(w.Width) {
		a.X = float64(w.Width)
		a.VX = -a.VX
	}
}

// func (w *World) Tick(dt time.Duration) {
// 	w.Grid.Clear()
// 	for i := range w.Agents {
// 		agent := &w.Agents[i]
// 		agent.Tick(dt, w.Width, w.Height)

// 		w.ApplyBoundaries(agent)
// 		w.Grid.Insert(agent.ID, float32(agent.X), float32(agent.Z))
// 	}
// }

func (w *World) AgentsTick(dt time.Duration) []agents.Agent {
	n := len(w.Agents)
	if n == 0 {
		return nil
	}

	results := make([]agents.Agent, n)
	changedFlags := make([]bool, n)
	var wg sync.WaitGroup

	wg.Add(n)

	for i := range n {
		go func(i int) {
			defer wg.Done()
			agent := &w.Agents[i]
			changed := agent.Tick(dt, w)
			results[i] = *agent
			if changed {
				changedFlags[i] = true
			}
		}(i)
	}

	wg.Wait()

	var changedAgents []agents.Agent
	for i := range n {
		a := results[i]
		w.Grid.Insert(a.ID, a.X, a.Z, a.Width+a.X, a.Height+a.Z, false)
		if changedFlags[i] {
			changedAgents = append(changedAgents, a)
		}
	}

	return changedAgents
}

func (w *World) GetWorldSeed() int64 {
	return w.Seed
}

func NewWorld(Seed int64, Width, Height float64) World {
	return World{
		Seed:   Seed,
		Width:  Width,
		Height: Height,
		Grid: spatialhash.SpatialHash{
			CellSize: 1,
			Cells:    make(map[int64]*spatialhash.Cell),
		},
		rng: rand.New(rand.NewSource(Seed)),
		// Agents: NewWorld().Agents,
	}
}
