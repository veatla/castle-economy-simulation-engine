package navgrid

import (
	"container/heap"
	worldQuery "veatla/simulator/src/world-query"
)

// PathNeighbors returns neighboring points (nil for raw PathPoint).
func (p *PathPoint) PathNeighbors(q worldQuery.WorldQuery) []Pather {
	return nil
}

// PathNeighborCost returns the cost to move to another point.
func (p *PathPoint) PathNeighborCost(to Pather, q worldQuery.WorldQuery) float64 {
	if point, ok := to.(*PathPoint); ok {
		return distance(p.X, p.Z, point.X, point.Z)
	}
	return 0
}

// PathEstimatedCost returns heuristic estimate.
func (p *PathPoint) PathEstimatedCost(to Pather) float64 {
	if point, ok := to.(*PathPoint); ok {
		return heuristic(p.X, p.Z, point.X, point.Z)
	}
	return 0
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

// Path finds a path using the generic Pather A* algorithm.
func Path(start, goal Pather, q worldQuery.WorldQuery) (pather []Pather, found bool, dist float64) {
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

func reconstructPath(node *Node) []PathPoint {
	path := make([]PathPoint, 0)
	current := node
	for current != nil {
		if current.Pather != nil {
			path = append(path, *current.Pather.(*PathPoint))
		}
		current = current.Parent
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}
