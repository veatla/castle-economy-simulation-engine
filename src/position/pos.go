package position

import (
	"math"
	"math/rand/v2"
)

type Position struct {
	X, Z          float64
	Width, Height int
}

func GetRandomPositionFromWorld(worldWidth int, worldHeight int, width int, height int) Position {
	return Position{
		X:      math.Round(rand.Float64() * float64(worldWidth)),
		Z:      math.Round(rand.Float64() * float64(worldHeight)),
		Width:  width,
		Height: height,
	}
}
