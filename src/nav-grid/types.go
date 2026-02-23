package navgrid

import worldQuery "veatla/simulator/src/world-query"

// PathPoint is a point in world space for pathfinding.
type PathPoint struct {
	X, Z float64
}

// Pather is the interface for A* pathfinding.
type Pather interface {
	PathNeighbors(q worldQuery.WorldQuery) []Pather
	PathNeighborCost(to Pather, q worldQuery.WorldQuery) float64
	PathEstimatedCost(to Pather) float64
}

// Node is used by the generic Pather-based Path().
type Node struct {
	Pather Pather
	Cost   float64
	Rank   float64
	Parent *Node
	Open   bool
	Closed bool
	index  int
}

// AStarNode is a node in the coordinate-based A* search.
type AStarNode struct {
	point  *PathPoint
	Cost   float64
	Rank   float64
	Parent *AStarNode
	Open   bool
	Closed bool
	index  int
}
