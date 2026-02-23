package agents

import (
	"log"
	"math"
	worldQuery "veatla/simulator/src/world-query"
	navgrid "veatla/simulator/src/nav-grid"
)

func (agent *Agent) detectStuck(q worldQuery.WorldQuery) {
	if agent.Wandering.wait >= 0 {
		ddx := agent.Wandering.X - agent.X
		ddz := agent.Wandering.Z - agent.Z
		dist2 := ddx*ddx + ddz*ddz
		const reachDist = 0.5

		if dist2 < reachDist*reachDist {
			agent.stuck.counter = 0
			agent.stuck.lastX = agent.X
			agent.stuck.lastZ = agent.Z
			return
		}
	}

	dx := agent.X - agent.stuck.lastX
	dz := agent.Z - agent.stuck.lastZ
	distance := math.Sqrt(dx*dx + dz*dz)

	const epsilon = 0.001
	const replanCooldown = 50

	if distance < epsilon {
		agent.stuck.counter++
		if agent.stuck.counter > agent.stuck.threshold {
			log.Printf("Agent %s is STUCK at position (%.2f, %.2f) for %d ticks with distance %.2f",
				agent.ID.String()[:8], agent.X, agent.Z, agent.stuck.counter,
				distance,
			)

			if agent.stuck.counter%replanCooldown == 0 {
				const obstacleOffset = 1.0
				if agent.Wandering.X != agent.X || agent.Wandering.Z != agent.Z {
					path, found, _ := navgrid.AStarPath(agent.X, agent.Z, agent.Wandering.X, agent.Wandering.Z, q, obstacleOffset)
					if found && len(path) > 0 {
						agent.path.path = path
						if len(path) > 1 {
							agent.path.pathIndex = 1
						} else {
							agent.path.pathIndex = 0
						}
					} else {
						agent.Wandering = agent.setWanderingTargetWithRadius(q, 8.0)
					}
				}
			}
			agent.stuck.counter = 0
			agent.stuck.lastX = agent.X
			agent.stuck.lastZ = agent.Z
		}
	} else {
		agent.stuck.counter = 0
		agent.stuck.lastX = agent.X
		agent.stuck.lastZ = agent.Z
	}
}
