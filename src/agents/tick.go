package agents

import (
	"math"
	"time"
	worldQuery "veatla/simulator/src/world-query"
)

func (agent *Agent) Tick(dt time.Duration, q worldQuery.WorldQuery) bool {
	oldX, oldZ := agent.X, agent.Z

	if agent.Wandering.wait <= 0 {
		agent.logWanderingEvent(q)
		agent.Wandering = agent.SetWanderingTarget(q)
	}

	dx := agent.Wandering.X - agent.X
	dz := agent.Wandering.Z - agent.Z
	dist2 := dx*dx + dz*dz
	const reachDist = 0.5

	if dist2 < reachDist*reachDist {
		agent.Wandering.wait -= dt
		return math.Abs(oldX-agent.X) > 1e-9 || math.Abs(oldZ-agent.Z) > 1e-9
	}

	agent.MoveTorwardsWanderingTarget(q)
	agent.detectStuck(q)

	return math.Abs(oldX-agent.X) > 1e-9 || math.Abs(oldZ-agent.Z) > 1e-9
}
