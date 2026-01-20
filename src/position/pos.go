package position

import (
	"math"
	"math/rand/v2"
)

type SolidCell struct {
	X int
	Z int
}

func GetRandomPositionFromWorld(worldWidth int, worldHeight int) SolidCell {
	return SolidCell{
		X: int(math.Round(rand.Float64() * float64(worldWidth))),
		Z: int(math.Round(rand.Float64() * float64(worldHeight))),
	}
}
