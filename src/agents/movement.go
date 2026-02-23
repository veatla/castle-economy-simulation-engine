package agents

import (
	"math"
	navgrid "veatla/simulator/src/nav-grid"
	"veatla/simulator/src/utils"
	worldQuery "veatla/simulator/src/world-query"
)

func (agent *Agent) MoveTorwardsWanderingTarget(q worldQuery.WorldQuery) {
	if agent.path.path != nil && agent.path.pathIndex < len(agent.path.path) {
		target := agent.path.path[agent.path.pathIndex]
		dx := target.X - agent.X
		dz := target.Z - agent.Z
		dist := math.Sqrt(dx*dx + dz*dz)
		if dist < 0.2 {
			agent.path.pathIndex++
			return
		}

		dx /= dist
		dz /= dist
		agent.VX = dx * (agent.baseSpeed + agent.Wandering.speed)
		agent.VZ = dz * (agent.baseSpeed + agent.Wandering.speed)
		agent.NoPath = false
		nextX := agent.X + agent.VX
		nextZ := agent.Z + agent.VZ

		if !q.IsPointBlocked(nextX, nextZ) {
			worldWidth, worldHeight := q.GetBoundaries()
			agent.X = utils.Clamp(nextX, 0, worldWidth)
			agent.Z = utils.Clamp(nextZ, 0, worldHeight)
			return
		}

		if !q.IsPointBlocked(target.X, target.Z) {
			worldWidth, worldHeight := q.GetBoundaries()
			agent.X = utils.Clamp(target.X, 0, worldWidth)
			agent.Z = utils.Clamp(target.Z, 0, worldHeight)
			agent.path.pathIndex++
			return
		}

		agent.navigateWithAStar(q)
		return
	}

	dx := agent.Wandering.X - agent.X
	dz := agent.Wandering.Z - agent.Z
	dist := math.Sqrt(dx*dx + dz*dz)
	if dist < 1e-6 {
		agent.VX = 0
		agent.VZ = 0
		return
	}
	dx /= dist
	dz /= dist
	step := agent.baseSpeed + agent.Wandering.speed
	nextX := agent.X + dx*step
	nextZ := agent.Z + dz*step
	if !q.IsPointBlocked(nextX, nextZ) {
		worldWidth, worldHeight := q.GetBoundaries()
		agent.X = utils.Clamp(nextX, 0, worldWidth)
		agent.Z = utils.Clamp(nextZ, 0, worldHeight)
		agent.VX = dx * step
		agent.VZ = dz * step
		return
	}
	agent.navigateWithAStar(q)
}

func (agent *Agent) navigateWithAStar(q worldQuery.WorldQuery) {
	const obstacleOffset = 1.0

	path, found, _ := navgrid.AStarPath(
		agent.X, agent.Z,
		agent.Wandering.X, agent.Wandering.Z,
		q, obstacleOffset,
	)

	if !found || len(path) == 0 {
		agent.Wandering = agent.SetWanderingTarget(q)
		agent.stuck.counter = 0
		agent.stuck.lastX = agent.X
		agent.stuck.lastZ = agent.Z
		return
	}

	agent.path.path = path
	if len(path) > 1 {
		agent.path.pathIndex = 1
	} else {
		agent.path.pathIndex = 0
	}
}
