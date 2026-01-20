package world

import (
	"example/hello/src/agents"
	"example/hello/src/position"
	"slices"

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
		agent := &w.Agents[i]

		var cell = position.SolidCell{
			X: int(agent.X),
			Z: int(agent.Z),
		}
		agent.X = agent.X + agent.VX*dt
		agent.Z = agent.X + agent.VZ*dt

		var newCell = position.SolidCell{
			X: int(agent.X),
			Z: int(agent.Z),
		}

		if newCell != cell {
			index := slices.Index(w.Grid[cell], agent.ID)
			if index != -1 {
				w.Grid[cell] = slices.Delete(w.Grid[cell], index, index+1)
			}
		}

		if agent.X < 0 || agent.X > float64(w.Width) {
			agent.VX = -agent.VX
		}

		if agent.Z < 0 || agent.Z > float64(w.Width) {
			agent.VZ = -agent.VZ
		}
	}
}

func NewWorld(Width int, Height int) World {
	return World{
		Width:  Width,
		Height: Height,
		// Agents: NewWorld().Agents,
	}
}
