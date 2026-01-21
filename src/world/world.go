package world

import (
	"example/hello/src/agents"
	"example/hello/src/position"

	"github.com/google/uuid"
)

type WorldID string

type World struct {
	Width  int
	Height int
	Agents []agents.Agent
	Grid   map[position.SolidCell][]uuid.UUID
}

func (w *World) Tick(dt float64) {
	for i := range w.Agents {
		w.Agents[i].MoveAgent(w.Grid, dt, w.Width, w.Height)
	}
}

func NewWorld(Width int, Height int) World {
	return World{
		Width:  Width,
		Height: Height,
		// Agents: NewWorld().Agents,
	}
}
