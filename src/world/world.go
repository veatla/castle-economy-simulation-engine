package world

import (
	"example/hello/src/agents"
	spatialhash "example/hello/src/spatial-hash"
	"time"
)

type WorldID string

type World struct {
	Width  int
	Height int
	Agents []agents.Agent
	Grid   spatialhash.SpatialHash
}

func (w *World) Tick(dt time.Duration) {
	w.Grid.Clear()
	for i := range w.Agents {
		agent := &w.Agents[i]
		agent.Tick(dt, w.Width, w.Height)

		w.ApplyBoundaries(agent)
		w.Grid.Insert(agent.ID, float32(agent.X), float32(agent.Z))
	}
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

	if a.Z < 0 {
		a.Z = 0
		a.VZ = -a.VZ
	}
	if a.Z > float64(w.Height) {
		a.Z = float64(w.Height)
		a.VZ = -a.VZ
	}
}
func NewWorld(Width int, Height int) World {
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
