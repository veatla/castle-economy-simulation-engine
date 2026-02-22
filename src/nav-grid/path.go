package navgrid

import (
	"container/heap"
	"math"
	worldQuery "veatla/simulator/src/world-query"
)

// PathPoint represents a point in world space for pathfinding
type PathPoint struct {
	X, Z float64
}

// PathNeighbors returns neighboring points
func (p *PathPoint) PathNeighbors(q worldQuery.WorldQuery) []Pather {
	return nil
}

// PathNeighborCost returns the cost to move to another point
func (p *PathPoint) PathNeighborCost(to Pather, q worldQuery.WorldQuery) float64 {
	if point, ok := to.(*PathPoint); ok {
		return distance(p.X, p.Z, point.X, point.Z)
	}
	return 0
}

// PathEstimatedCost returns heuristic estimate
func (p *PathPoint) PathEstimatedCost(to Pather) float64 {
	if point, ok := to.(*PathPoint); ok {
		return heuristic(p.X, p.Z, point.X, point.Z)
	}
	return 0
}

// Pather interface for A* pathfinding
type Pather interface {
	PathNeighbors(q worldQuery.WorldQuery) []Pather
	PathNeighborCost(to Pather, q worldQuery.WorldQuery) float64
	PathEstimatedCost(to Pather) float64
}

type Node struct {
	Pather Pather
	Cost   float64
	Rank   float64
	Parent *Node
	Open   bool
	Closed bool
	index  int
}

type NodeMap map[Pather]*Node

func (nm NodeMap) Get(pather Pather) *Node {
	n, ok := nm[pather]
	if !ok {
		n = &Node{Pather: pather}
		nm[pather] = n
	}

	return n
}

// Path finds a path using A* algorithm
func Path(start, goal Pather, q worldQuery.WorldQuery) (pather []Pather, found bool, distance float64) {
	nm := make(NodeMap)
	nq := &PriorityQueue{}

	heap.Init(nq)

	fromNode := nm.Get(start)
	fromNode.Open = true
	fromNode.Cost = 0
	heap.Push(nq, fromNode)

	for {
		if nq.Len() == 0 {
			return nil, false, 0.0
		}

		current := heap.Pop(nq).(*Node)
		current.Open = false
		current.Closed = true

		if current == nm.Get(goal) {
			p := []Pather{}
			curr := current
			for curr != nil {
				p = append(p, curr.Pather)
				curr = curr.Parent
			}
			return p, true, current.Cost
		}

		for _, neighbor := range current.Pather.PathNeighbors(q) {
			cost := current.Cost + current.Pather.PathNeighborCost(neighbor, q)
			neighborNode := nm.Get(neighbor)

			if cost < neighborNode.Cost || (!neighborNode.Open && !neighborNode.Closed) {
				if neighborNode.Open {
					heap.Remove(nq, neighborNode.index)
				}
				neighborNode.Cost = cost
				neighborNode.Open = true
				neighborNode.Closed = false
				neighborNode.Rank = cost + neighbor.PathEstimatedCost(goal)
				neighborNode.Parent = current
				heap.Push(nq, neighborNode)
			}
		}
	}
}

// AStarPath finds path between two world coordinates with obstacle avoidance and offset
func AStarPath(startX, startZ, goalX, goalZ float64, q worldQuery.WorldQuery, obstacleOffset float64) ([]PathPoint, bool, float64) {
	// Use internal search to get raw path, then validate/repair
	const cellSize = 2.0 // balanced: smaller than 3 but larger than 1 for efficiency
	rawPath, found, total := aStarSearch(startX, startZ, goalX, goalZ, q, obstacleOffset, cellSize)
	if !found || len(rawPath) == 0 {
		return nil, false, 0.0
	}

	// Validate path (no repair - just reject if blocked)
	if !validatePathSimple(rawPath, q, obstacleOffset, cellSize) {
		return nil, false, 0.0
	}

	// Simplify path by removing redundant waypoints
	// simplified := simplifyPath(rawPath)
	return rawPath, true, total
}

// aStarSearch performs the A* search but does not run post-validation. It returns the raw path.
func aStarSearch(startX, startZ, goalX, goalZ float64, q worldQuery.WorldQuery, obstacleOffset, cellSize float64) ([]PathPoint, bool, float64) {
	const maxPath = 1000

	// Don't block if goal is directly blocked, just check that it's within bounds
	// The path might not reach exactly the goal but can get close

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
	var bestNode *AStarNode // track best node reached in case we can't reach exact goal

	for pq.Len() > 0 && pathCount < maxPath {
		pathCount++
		current := heap.Pop(pq).(*AStarNode)
		currentKey := pointKey(current.point.X, current.point.Z)

		delete(openSet, currentKey)
		closedSet[currentKey] = current

		// Track best approach to goal
		if bestNode == nil || heuristic(current.point.X, current.point.Z, goalX, goalZ) < heuristic(bestNode.point.X, bestNode.point.Z, goalX, goalZ) {
			bestNode = current
		}

		// Check if we reached goal (within cellSize tolerance)
		if distance(current.point.X, current.point.Z, goalX, goalZ) < cellSize {
			path := reconstructAStarPath(current)
			totalDist := current.Cost
			return path, true, totalDist
		}

		// Generate neighbors
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

	// If we couldn't reach exact goal but got reasonably close, return best approach
	if bestNode != nil && heuristic(bestNode.point.X, bestNode.point.Z, goalX, goalZ) < cellSize*3 {
		path := reconstructAStarPath(bestNode)
		totalDist := bestNode.Cost
		return path, true, totalDist
	}

	return nil, false, 0.0
}

// AStarNode represents a node in the A* search
type AStarNode struct {
	point  *PathPoint
	Cost   float64
	Rank   float64
	Parent *AStarNode
	Open   bool
	Closed bool
	index  int
}

func pointKey(x, z float64) string {
	return string(rune(int(x*1000))) + string(rune(int(z*1000)))
}

func heuristic(x1, z1, x2, z2 float64) float64 {
	// Euclidean distance heuristic
	dx := x2 - x1
	dz := z2 - z1
	return math.Sqrt(dx*dx + dz*dz)
}

func distance(x1, z1, x2, z2 float64) float64 {
	dx := x2 - x1
	dz := z2 - z1
	return math.Sqrt(dx*dx + dz*dz)
}

// getNeighbors returns valid neighbor points with obstacle offset
func getNeighbors(current *PathPoint, cellSize, obstacleOffset float64, q worldQuery.WorldQuery) []*PathPoint {
	neighbors := make([]*PathPoint, 0, 8)

	// 8-directional movement
	directions := []struct{ dx, dz float64 }{
		{cellSize, 0},
		{-cellSize, 0},
		{0, cellSize},
		{0, -cellSize},
		{cellSize * 0.707, cellSize * 0.707},   // diagonal
		{-cellSize * 0.707, cellSize * 0.707},  // diagonal
		{cellSize * 0.707, -cellSize * 0.707},  // diagonal
		{-cellSize * 0.707, -cellSize * 0.707}, // diagonal
	}

	for _, dir := range directions {
		newX := current.X + dir.dx
		newZ := current.Z + dir.dz

		// Check if point is blocked
		if q.IsPointBlocked(newX, newZ) {
			continue
		}

		// Verify with offset (maintain distance from obstacles)
		// Be more lenient: if we can't maintain offset, try without it
		if !isPointValidWithOffset(newX, newZ, obstacleOffset, q) {
			// Try with reduced offset as fallback
			if !isPointValidWithOffset(newX, newZ, obstacleOffset*0.5, q) {
				continue
			}
		}

		neighbors = append(neighbors, &PathPoint{X: newX, Z: newZ})
	}

	return neighbors
}

// isPointValidWithOffset checks if a point is valid and maintains offset from obstacles
func isPointValidWithOffset(x, z, offset float64, q worldQuery.WorldQuery) bool {
	if q.IsPointBlocked(x, z) {
		return false
	}

	// Check points around in a circle to ensure offset distance
	checkPoints := 8
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

func reconstructPath(node *Node) []PathPoint {
	path := make([]PathPoint, 0)
	current := node
	for current != nil {
		if current.Pather != nil {
			path = append(path, *current.Pather.(*PathPoint))
		}
		current = current.Parent
	}
	// Reverse path to get start -> goal order
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
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
	// Reverse path to get start -> goal order
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

// validatePathSimple validates path without trying to repair (simple and fast)
func validatePathSimple(path []PathPoint, q worldQuery.WorldQuery, obstacleOffset, cellSize float64) bool {
	if len(path) == 0 {
		return false
	}

	// Check each waypoint for basic validity (not blocked, no hard obstacles)
	for _, p := range path {
		if q.IsPointBlocked(p.X, p.Z) {
			return false
		}
	}

	// Quick segment check: ensure waypoints don't have obstacles between them
	for i := 0; i < len(path)-1; i++ {
		a := path[i]
		b := path[i+1]
		dx := b.X - a.X
		dz := b.Z - a.Z
		segLen := math.Sqrt(dx*dx + dz*dz)
		if segLen == 0 {
			continue
		}
		// Sample only a few points along segment for speed
		for s := 1; s <= 2; s++ {
			t := float64(s) / 3.0
			sx := a.X + dx*t
			sz := a.Z + dz*t
			// Only check if blocked, don't enforce strict offset during validation
			// (paths computed with offset will naturally have offset)
			if q.IsPointBlocked(sx, sz) {
				return false
			}
		}
	}

	return true
}

// attemptRepair tries to replan from path[idx] to the original goal using raw aStarSearch.
// If successful, it returns a new spliced path that replaces path[idx+1:] with the repaired subpath.
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

	// Ensure no duplicate start point
	if len(newPath) > 0 {
		if math.Abs(newPath[0].X-start.X) < 1e-6 && math.Abs(newPath[0].Z-start.Z) < 1e-6 {
			// splice: keep path[:idx+1], then append newPath[1:]
			spliced := make([]PathPoint, 0, idx+1+len(newPath)-1)
			spliced = append(spliced, path[:idx+1]...)
			if len(newPath) > 1 {
				spliced = append(spliced, newPath[1:]...)
			}
			return spliced, true
		}
	}

	// If newPath doesn't start at start, conservatively prepend start then the new path
	spliced := make([]PathPoint, 0, idx+1+len(newPath))
	spliced = append(spliced, path[:idx+1]...)
	spliced = append(spliced, newPath...)
	return spliced, true
}

// func (Pather *Pather) MovementCost(current, destination Point, q worldQuery.WorldQuery) float64 {
// 	if q.IsPointBlocked(destination.X, destination.Z) {
// 		return -1.0
// 	}

// 	dx := destination.X - current.X
// 	dz := destination.Z - current.Z

// 	return (dx*dx + dz*dz)
// }

// func (Pather *Pather) ManhattanDistance(a, b Point) float64 {
// 	dx := b.X - a.X
// 	dz := b.Z - a.Z

// 	if dx < 0 {
// 		dx = -dx
// 	}

// 	if dz < 0 {
// 		dz = -dz
// 	}

// 	return dx + dz
// }

// simplifyPath removes redundant waypoints that lie on a straight line
func simplifyPath(path []PathPoint) []PathPoint {
	if len(path) <= 2 {
		return path
	}

	simplified := []PathPoint{path[0]}
	for i := 1; i < len(path)-1; i++ {
		prev := simplified[len(simplified)-1]
		curr := path[i]
		next := path[i+1]

		// Check if curr is roughly collinear with prev and next
		// Using cross product: if ~0, points are collinear
		dx1 := curr.X - prev.X
		dz1 := curr.Z - prev.Z
		dx2 := next.X - curr.X
		dz2 := next.Z - curr.Z

		crossProduct := math.Abs(dx1*dz2 - dz1*dx2)
		const collinearThreshold = 0.1 // tolerance for collinearity

		if crossProduct > collinearThreshold {
			// Not collinear, keep this waypoint
			simplified = append(simplified, curr)
		}
	}
	// Always include the final waypoint
	simplified = append(simplified, path[len(path)-1])

	return simplified
}

// func (Pather *Pather) Heuristic(a, b Point) float64 {
// 	return Pather.ManhattanDistance(a, b)
// }
