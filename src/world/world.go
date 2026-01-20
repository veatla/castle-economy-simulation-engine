package world

import (
	"example/hello/src/agents"
	"slices"

	"github.com/google/uuid"
)

type WorldID string
type Cell struct {
	X float64
	Z float64
}
type World struct {
	Width  int
	Height int
	Agents []agents.Agent
	Grid   map[Cell][]uuid.UUID
}

func (w *World) Tick(dt float64) {
	for i := range w.Agents {
		agent := &w.Agents[i]

		var cell = Cell{
			X: agent.X,
			Z: agent.Z,
		}
		agent.X = agent.X + agent.VX*dt
		agent.Z = agent.X + agent.VZ*dt

		var newCell = Cell{
			X: agent.X,
			Z: agent.Z,
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

		// fmt.Printf("Agent #%v on X: %f Z: %f \n", agent.ID, agent.X, agent.Z)
	}
}

func NewWorld(Width int, Height int) World {
	return World{
		Width:  Width,
		Height: Height,
		// Agents: NewWorld().Agents,
	}
}
