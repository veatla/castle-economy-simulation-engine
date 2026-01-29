package worldQuery

type WorldQuery interface {
	IsPointBlocked(x, z float64) bool
	RandomFloat() float64
	GetWorldSeed() int64
	GetBoundaries() (width, height float64)
}
