package navgrid

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Rank < pq[j].Rank
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]

	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	no := x.(*Node)
	no.index = n
	*pq = append(*pq, no)
}
func (pq *PriorityQueue) Pop() interface{} {

	old := *pq
	n := len(old)
	no := old[n-1]
	no.index = -1
	*pq = old[0 : n-1]
	return no
}

// AStarPriorityQueue for A* node pathfinding
type AStarPriorityQueue []*AStarNode

func (pq AStarPriorityQueue) Len() int {
	return len(pq)
}

func (pq AStarPriorityQueue) Less(i, j int) bool {
	return pq[i].Rank < pq[j].Rank
}

func (pq AStarPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *AStarPriorityQueue) Push(x interface{}) {
	n := len(*pq)
	no := x.(*AStarNode)
	no.index = n
	*pq = append(*pq, no)
}

func (pq *AStarPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	no := old[n-1]
	no.index = -1
	*pq = old[0 : n-1]
	return no
}
