package spatialhash

type Cell struct {
	agents []int
}
type SpatialHash struct {
	cellSize float32
	cells    map[int64]*Cell
}

func hashCell(x, z int32) int64 {
	return int64(x)<<32 | int64(z)&0xffffffff
}

func (s *SpatialHash) cellFor(x, z float32) (int32, int32) {
	return int32(x / s.cellSize), int32(z / s.cellSize)
}
func (s *SpatialHash) Clear() {
	for _, c := range s.cells {
		c.agents = c.agents[:0]
	}

}

func (s *SpatialHash) Insert(agentID int, x, z float32) {
	cx, cz := s.cellFor(x, z)
	key := hashCell(cx, cz)

	cell, ok := s.cells[key]
	if !ok {
		cell = &Cell{}
		s.cells[key] = cell
	}

	cell.agents = append(cell.agents, agentID)
}

func (s *SpatialHash) Nearby(x, z float32) []int {
	cx, cz := s.cellFor(x, z)

	var result []int

	for dx := -1; dx <= 1; dx++ {
		for dz := -1; dz <= 1; dz++ {
			key := hashCell(cx+int32(dx), cz+int32(dz))
			if cell := s.cells[key]; cell != nil {
				result = append(result, cell.agents...)
			}
		}
	}

	return result
}
