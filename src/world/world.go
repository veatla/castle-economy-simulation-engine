package world

import (
	"example/hello/src/agents"
	spatialhash "example/hello/src/spatial-hash"
	"sync"
	"time"
)

type WorldID string

type World struct {
	Width  float64
	Height float64
	Agents []agents.Agent
	Grid   spatialhash.SpatialHash
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

func (w *World) Tick(dt time.Duration) []agents.Agent {
	w.Grid.Clear()
	n := len(w.Agents)
	if n == 0 {
		return nil
	}

	results := make([]agents.Agent, n)
	changedFlags := make([]bool, n)
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			agent := &w.Agents[i]
			changed := agent.Tick(dt, w.Width, w.Height)
			results[i] = *agent
			if changed {
				changedFlags[i] = true
			}
		}(i)
	}

	wg.Wait()

	var changedAgents []agents.Agent
	for i := 0; i < n; i++ {
		a := results[i]
		w.Grid.Insert(a.ID, float32(a.X), float32(a.Z))
		if changedFlags[i] {
			changedAgents = append(changedAgents, a)
		}
	}

	return changedAgents
}

func NewWorld(Width, Height float64) World {
	return World{
		Width:  Width,
		Height: Height,
		Grid: spatialhash.SpatialHash{
			CellSize: 100,
			Cells:    make(map[int64]*spatialhash.Cell),
		},
		// Agents: NewWorld().Agents,
	}
}
