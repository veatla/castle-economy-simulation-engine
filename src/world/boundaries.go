package world

import "veatla/simulator/src/agents"

func (w *World) ApplyBoundaries(a *agents.Agent) {
	if a.X < 0 {
		a.X = 0
		a.VX = -a.VX
	}
	if a.X > w.Width {
		a.X = w.Width
		a.VX = -a.VX
	}
}
