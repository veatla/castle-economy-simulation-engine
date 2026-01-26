package utils

import "math"

func Clamp(value, min, max float64) float64 {
	return math.Min(math.Max(value, min), max)
}
