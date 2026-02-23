package world

import (
	"sync"
	"time"

	"veatla/simulator/src/agents"
)

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
