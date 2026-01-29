package navgrid

type Path struct {
	Points []Point
}

type Point struct {
	X, Z float64
}

func FindPath(grid *NavGrid, startX, startZ, endX, endZ float64) []CellPos {
	return []CellPos{
		{
			X: int(startX / grid.CellSize),
			Z: int(startZ / grid.CellSize),
		},
		{
			X: int(endX / grid.CellSize),
			Z: int(endZ / grid.CellSize),
		},
	}
}

type CellPos struct {
	X, Z int
}
