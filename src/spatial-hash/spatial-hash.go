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
func (s *SpatialHash) Clear(includeStructures bool) {
	for _, c := range s.Cells {
		c.agents = c.agents[:0]

		if includeStructures == true {
			c.obstacles = c.obstacles[:0]
		}
	}
}

func (s *SpatialHash) Insert(id uuid.UUID, x, z float64, x2, z2 float64, structure bool) {
	for dx := x; dx <= x2; dx += float64(s.CellSize) {
		for dz := z; dz <= z2; dz += float64(s.CellSize) {
			s.insertAtCell(id, dx, dz, structure)
		}
	}
}

func (s *SpatialHash) insertAtCell(id uuid.UUID, x, z float64, structure bool) {
	cx, cz := s.cellFor(float32(x), float32(z))
	key := HashCell(cx, cz)

	cell, ok := s.Cells[key]
	if !ok {
		cell = &Cell{}
		s.Cells[key] = cell
	}

	if structure == false {
		cell.agents = append(cell.agents, id)
	} else {
		cell.obstacles = append(cell.obstacles, id)
	}
}

func (s *SpatialHash) Nearby(x, z float32, structure bool) []uuid.UUID {
	cx, cz := s.cellFor(x, z)

	var result []uuid.UUID

	for dx := -1; dx <= 1; dx++ {
		for dz := -1; dz <= 1; dz++ {
			key := HashCell(cx+int32(dx), cz+int32(dz))
			if cell := s.Cells[key]; cell != nil {
				if structure == false {
					result = append(result, cell.agents...)
				} else {
					result = append(result, cell.obstacles...)
				}
			}
		}
	}

	return result
}

func (s *SpatialHash) IsPointBlocked(x, z float64) bool {
	cx, cz := s.cellFor(float32(x), float32(z))
	key := HashCell(cx, cz)

	cell, ok := s.Cells[key]
	if !ok {
		return false
	}

	return len(cell.obstacles) > 0
}
