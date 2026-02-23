package navgrid

import (
	"math"
	worldQuery "veatla/simulator/src/world-query"
)

func validatePathSimple(path []PathPoint, q worldQuery.WorldQuery, obstacleOffset, cellSize float64) bool {
	if len(path) == 0 {
		return false
	}

	for _, p := range path {
		if q.IsPointBlocked(p.X, p.Z) {
			return false
		}
	}

	for i := 0; i < len(path)-1; i++ {
		a := path[i]
		b := path[i+1]
		dx := b.X - a.X
		dz := b.Z - a.Z
		for s := 1; s <= 2; s++ {
			t := float64(s) / 3.0
			sx := a.X + dx*t
			sz := a.Z + dz*t
			if q.IsPointBlocked(sx, sz) {
				return false
			}
		}
	}
	return true
}

func attemptRepair(path []PathPoint, idx int, q worldQuery.WorldQuery, obstacleOffset, cellSize float64) ([]PathPoint, bool) {
	if idx < 0 || idx >= len(path) {
		return nil, false
	}
	start := path[idx]
	goal := path[len(path)-1]

	newPath, found, _ := aStarSearch(start.X, start.Z, goal.X, goal.Z, q, obstacleOffset, cellSize)
	if !found || len(newPath) == 0 {
		return nil, false
	}

	if len(newPath) > 0 {
		if math.Abs(newPath[0].X-start.X) < 1e-6 && math.Abs(newPath[0].Z-start.Z) < 1e-6 {
			spliced := make([]PathPoint, 0, idx+1+len(newPath)-1)
			spliced = append(spliced, path[:idx+1]...)
			if len(newPath) > 1 {
				spliced = append(spliced, newPath[1:]...)
			}
			return spliced, true
		}
	}

	spliced := make([]PathPoint, 0, idx+1+len(newPath))
	spliced = append(spliced, path[:idx+1]...)
	spliced = append(spliced, newPath...)
	return spliced, true
}

func simplifyPath(path []PathPoint) []PathPoint {
	if len(path) <= 2 {
		return path
	}

	simplified := []PathPoint{path[0]}
	for i := 1; i < len(path)-1; i++ {
		prev := simplified[len(simplified)-1]
		curr := path[i]
		next := path[i+1]

		dx1 := curr.X - prev.X
		dz1 := curr.Z - prev.Z
		dx2 := next.X - curr.X
		dz2 := next.Z - curr.Z

		crossProduct := math.Abs(dx1*dz2 - dz1*dx2)
		const collinearThreshold = 0.1

		if crossProduct > collinearThreshold {
			simplified = append(simplified, curr)
		}
	}
	simplified = append(simplified, path[len(path)-1])
	return simplified
}
