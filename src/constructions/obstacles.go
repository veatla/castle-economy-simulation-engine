package constructions

import "github.com/google/uuid"

type Obstacle struct {
	ID         uuid.UUID
	MinX, MinZ float64
	MaxX, MaxZ float64
}

func CreateObstacle(MinX, MinZ, MaxX, MaxZ float64) Obstacle {
	return Obstacle{

		MinX: MinX,
		MinZ: MinZ,
		MaxX: MaxX,
		MaxZ: MaxZ,
	}
}
