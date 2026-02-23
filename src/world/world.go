package world

import (
	"math/rand"

	spatialhash "veatla/simulator/src/spatial-hash"
)

func NewWorld(seed int64, width, height float64) World {
	return World{
		Seed:   seed,
		Width:  width,
		Height: height,
		Grid: spatialhash.SpatialHash{
			CellSize: 1,
			Cells:    make(map[int64]*spatialhash.Cell),
		},
		rng: rand.New(rand.NewSource(seed)),
	}
}
