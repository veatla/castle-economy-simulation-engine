package position

import "math/rand/v2"

type Position struct {
	X, Z          float32
	Width, Height float32
}

func GetRandomPositionFromWorld(worldWidth float32, worldHeight float32, width float32, height float32) Position {
	return Position{
		X:      rand.Float32() * worldWidth,
		Z:      rand.Float32() * worldHeight,
		Width:  width,
		Height: height,
	}
}
