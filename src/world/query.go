package world

func (w *World) IsPointBlocked(x, z float64) bool {
	return w.Grid.IsPointBlocked(x, z)
}

func (w *World) RandomFloat() float64 {
	return w.rng.Float64()
}

func (w *World) GetBoundaries() (width, height float64) {
	return w.Width, w.Height
}

func (w *World) GetWorldSeed() int64 {
	return w.Seed
}
