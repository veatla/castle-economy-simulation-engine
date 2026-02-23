package navgrid

import "math"

func heuristic(x1, z1, x2, z2 float64) float64 {
	dx := x2 - x1
	dz := z2 - z1
	return math.Sqrt(dx*dx + dz*dz)
}

func distance(x1, z1, x2, z2 float64) float64 {
	dx := x2 - x1
	dz := z2 - z1
	return math.Sqrt(dx*dx + dz*dz)
}

func pointKey(x, z float64) string {
	return string(rune(int(x*1000))) + string(rune(int(z*1000)))
}
