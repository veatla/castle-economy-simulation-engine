package world

import (
	"example/hello/src/agents"
	spatialhash "example/hello/src/spatial-hash"
)

type WorldID string

type World struct {
	Width  int
	Height int
	Agents []agents.Agent
	Grid   spatialhash.SpatialHash
}

func (w *World) Tick(dt float64) {
	w.Grid.Clear()
	for i := range w.Agents {
		agent := w.Agents[i]
		agent.Tick()
		w.Grid.Insert(agent.ID, float32(agent.X), float32(agent.Z))
	}
}

func NewWorld(Width int, Height int) World {
	return World{
		Width:  Width,
		Height: Height,
		// Agents: NewWorld().Agents,
	}
}
