package world

import (
	"example/hello/src/agents"
)

type WorldID string

type World struct {
	Width  float32
	Height float32
	Agents []agents.Agent
}

func (w *World) Tick(dt float32) {
	for i := range w.Agents {
		agent := &w.Agents[i]

		agent.X += agent.VX * dt
		agent.Z += agent.VZ * dt

		if agent.X < 0 || agent.X > w.Width {
			agent.VX = -agent.VX
		}
		if agent.Z < 0 || agent.Z > w.Width {
			agent.VZ = -agent.VZ
		}
		// fmt.Printf("Agent #%v on X: %f Z: %f \n", agent.ID, agent.X, agent.Z)
	}
}

func NewWorld() World {
	return World{
		Width:  100,
		Height: 100,
		// Agents: NewWorld().Agents,
	}
}
