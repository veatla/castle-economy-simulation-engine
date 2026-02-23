package agents

import (
	"math"
	"time"
	navgrid "veatla/simulator/src/nav-grid"
	"veatla/simulator/src/utils"
	worldQuery "veatla/simulator/src/world-query"
)

func (agent *Agent) SetWanderingTarget(q worldQuery.WorldQuery) Wandering {
	return agent.setWanderingTargetWithRadius(q, 30.0)
}

func (agent *Agent) setWanderingTargetWithRadius(q worldQuery.WorldQuery, maxRadius float64) Wandering {
	const obstacleOffset = 1.0

	angle := agent.rng.Float64() * 2 * math.Pi
	radius := math.Sqrt(agent.rng.Float64()) * maxRadius

	dx := math.Cos(angle) * radius
	dz := math.Sin(angle) * radius

	tx := agent.X + dx
	tz := agent.Z + dz

	worldWidth, worldHeight := q.GetBoundaries()

	for tx < 0 || tx > worldWidth || tz < 0 || tz > worldHeight || q.IsPointBlocked(tx, tz) {
		angle = agent.rng.Float64() * 2 * math.Pi
		radius = math.Sqrt(agent.rng.Float64()) * maxRadius
		dx = math.Cos(angle) * radius
		dz = math.Sin(angle) * radius
		tx = agent.X + dx
		tz = agent.Z + dz
	}

	var targetX, targetZ float64
	if offsetX, offsetZ, ok := utils.FindOffsetPosition(tx, tz, obstacleOffset, q); ok {
		targetX = offsetX
		targetZ = offsetZ
	} else {
		targetX = utils.Clamp(tx, 0, worldWidth)
		targetZ = utils.Clamp(tz, 0, worldHeight)
	}

	if path, found, _ := navgrid.AStarPath(agent.X, agent.Z, targetX, targetZ, q, obstacleOffset); found && len(path) > 0 {
		agent.path.path = path
		agent.NoPath = false
		if len(path) > 1 {
			agent.path.pathIndex = 1
		} else {
			agent.path.pathIndex = 0
		}
	} else {
		agent.path.path = nil
		agent.NoPath = true
		agent.path.pathIndex = 0
	}

	return Wandering{
		speed: 0.03 + agent.rng.Float64()*0.02,
		X:     targetX,
		Z:     targetZ,
		wait:  500 * time.Millisecond,
	}
}
