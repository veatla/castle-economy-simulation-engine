package navgrid

import (
	"container/heap"
	"math"
	worldQuery "veatla/simulator/src/world-query"
)

const maxPathLen = 1000

// AStarPath finds a path between two world coordinates with obstacle avoidance and offset.
func AStarPath(startX, startZ, goalX, goalZ float64, q worldQuery.WorldQuery, obstacleOffset float64) ([]PathPoint, bool, float64) {
	const cellSize = 2.0
	rawPath, found, total := aStarSearch(startX, startZ, goalX, goalZ, q, obstacleOffset, cellSize)
	if !found || len(rawPath) == 0 {
		return nil, false, 0.0
	}
	if !validatePathSimple(rawPath, q, obstacleOffset, cellSize) {
		return nil, false, 0.0
	}
	return rawPath, true, total
}

func aStarSearch(startX, startZ, goalX, goalZ float64, q worldQuery.WorldQuery, obstacleOffset, cellSize float64) ([]PathPoint, bool, float64) {
	openSet := make(map[string]*AStarNode)
	closedSet := make(map[string]*AStarNode)
	pq := &AStarPriorityQueue{}
	heap.Init(pq)

	startKey := pointKey(startX, startZ)
	startNode := &AStarNode{
		point: &PathPoint{X: startX, Z: startZ},
		Cost:  0,
		Rank:  heuristic(startX, startZ, goalX, goalZ),
	}
	openSet[startKey] = startNode
	heap.Push(pq, startNode)

	pathCount := 0
	var bestNode *AStarNode

	for pq.Len() > 0 && pathCount < maxPathLen {
		pathCount++
		current := heap.Pop(pq).(*AStarNode)
		currentKey := pointKey(current.point.X, current.point.Z)

		delete(openSet, currentKey)
		closedSet[currentKey] = current

		if bestNode == nil || heuristic(current.point.X, current.point.Z, goalX, goalZ) < heuristic(bestNode.point.X, bestNode.point.Z, goalX, goalZ) {
			bestNode = current
		}

		if distance(current.point.X, current.point.Z, goalX, goalZ) < cellSize {
			return reconstructAStarPath(current), true, current.Cost
		}

		for _, neighbor := range getNeighbors(current.point, cellSize, obstacleOffset, q) {
			neighborKey := pointKey(neighbor.X, neighbor.Z)

			if _, inClosed := closedSet[neighborKey]; inClosed {
				continue
			}

			moveCost := distance(current.point.X, current.point.Z, neighbor.X, neighbor.Z)
			newCost := current.Cost + moveCost

			if existingNode, inOpen := openSet[neighborKey]; inOpen {
				if newCost < existingNode.Cost {
					existingNode.Cost = newCost
					existingNode.Rank = newCost + heuristic(neighbor.X, neighbor.Z, goalX, goalZ)
					existingNode.Parent = current
					heap.Fix(pq, existingNode.index)
				}
			} else {
				neighborNode := &AStarNode{
					point:  neighbor,
					Cost:   newCost,
					Rank:   newCost + heuristic(neighbor.X, neighbor.Z, goalX, goalZ),
					Parent: current,
					Open:   true,
				}
				openSet[neighborKey] = neighborNode
				heap.Push(pq, neighborNode)
			}
		}
	}

	if bestNode != nil && heuristic(bestNode.point.X, bestNode.point.Z, goalX, goalZ) < cellSize*3 {
		return reconstructAStarPath(bestNode), true, bestNode.Cost
	}

	return nil, false, 0.0
}

func getNeighbors(current *PathPoint, cellSize, obstacleOffset float64, q worldQuery.WorldQuery) []*PathPoint {
	neighbors := make([]*PathPoint, 0, 8)
	directions := []struct{ dx, dz float64 }{
		{cellSize, 0},
		{-cellSize, 0},
		{0, cellSize},
		{0, -cellSize},
		{cellSize * 0.707, cellSize * 0.707},
		{-cellSize * 0.707, cellSize * 0.707},
		{cellSize * 0.707, -cellSize * 0.707},
		{-cellSize * 0.707, -cellSize * 0.707},
	}

	for _, dir := range directions {
		newX := current.X + dir.dx
		newZ := current.Z + dir.dz

		if q.IsPointBlocked(newX, newZ) {
			continue
		}
		if !isPointValidWithOffset(newX, newZ, obstacleOffset, q) {
			if !isPointValidWithOffset(newX, newZ, obstacleOffset*0.5, q) {
				continue
			}
		}
		neighbors = append(neighbors, &PathPoint{X: newX, Z: newZ})
	}
	return neighbors
}

func isPointValidWithOffset(x, z, offset float64, q worldQuery.WorldQuery) bool {
	if q.IsPointBlocked(x, z) {
		return false
	}
	const checkPoints = 8
	for i := 0; i < checkPoints; i++ {
		angle := float64(i) * 2 * math.Pi / float64(checkPoints)
		checkX := x + math.Cos(angle)*offset
		checkZ := z + math.Sin(angle)*offset
		if q.IsPointBlocked(checkX, checkZ) {
			return false
		}
	}
	return true
}

func reconstructAStarPath(node *AStarNode) []PathPoint {
	path := make([]PathPoint, 0)
	current := node
	for current != nil {
		if current.point != nil {
			path = append(path, *current.point)
		}
		current = current.Parent
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}
