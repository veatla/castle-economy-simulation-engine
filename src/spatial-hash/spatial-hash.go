package spatialhash

import "github.com/google/uuid"

type Cell struct {
	agents    []uuid.UUID
	obstacles []uuid.UUID
}
type SpatialHash struct {
	CellSize float32
	Cells    map[int64]*Cell
}

func HashCell(x, z int32) int64 {
	return int64(x)<<32 | int64(z)&0xffffffff
}

func (s *SpatialHash) cellFor(x, z float32) (int32, int32) {
	return int32(x / s.CellSize), int32(z / s.CellSize)
}
func (s *SpatialHash) Clear() {
	for _, c := range s.Cells {
		c.agents = c.agents[:0]
	}
}

func (s *SpatialHash) Insert(agentID uuid.UUID, x, z float32) {
	cx, cz := s.cellFor(x, z)
	key := HashCell(cx, cz)

	cell, ok := s.Cells[key]
	if !ok {
		cell = &Cell{}
		s.Cells[key] = cell
	}

	cell.agents = append(cell.agents, agentID)
}

func (s *SpatialHash) Nearby(x, z float32) []uuid.UUID {
	cx, cz := s.cellFor(x, z)

	var result []uuid.UUID

	for dx := -1; dx <= 1; dx++ {
		for dz := -1; dz <= 1; dz++ {
			key := HashCell(cx+int32(dx), cz+int32(dz))
			if cell := s.Cells[key]; cell != nil {
				result = append(result, cell.agents...)
			}
		}
	}

	return result
}
