package utils

import (
	"math"
	worldQuery "veatla/simulator/src/world-query"
)

// FindOffsetPosition finds a valid position offset from obstacles
func FindOffsetPosition(x, z, offsetDistance float64, q worldQuery.WorldQuery) (float64, float64, bool) {
	const numAngles = 16
	const step = 0.1

	// Try to find nearest valid position with offset
	for attempt := 0; attempt < 5; attempt++ {
		currentDist := offsetDistance + (float64(attempt) * step)

		for i := 0; i < numAngles; i++ {
			angle := float64(i) * 2 * math.Pi / float64(numAngles)
			offsetX := x + math.Cos(angle)*currentDist
			offsetZ := z + math.Sin(angle)*currentDist

			if !q.IsPointBlocked(offsetX, offsetZ) && isOffsetValid(offsetX, offsetZ, offsetDistance, q) {
				return offsetX, offsetZ, true
			}
		}
	}

	return x, z, false
}

// isOffsetValid checks if a position is valid and maintains offset from obstacles
func isOffsetValid(x, z, offsetDistance float64, q worldQuery.WorldQuery) bool {
	if q.IsPointBlocked(x, z) {
		return false
	}

	// Check 8 directions around the point
	const checkAngles = 8
	for i := 0; i < checkAngles; i++ {
		angle := float64(i) * 2 * math.Pi / float64(checkAngles)
		checkX := x + math.Cos(angle)*offsetDistance
		checkZ := z + math.Sin(angle)*offsetDistance

		if q.IsPointBlocked(checkX, checkZ) {
			return false
		}
	}

	return true
}

// GetSafeWanderingTarget tries to find a safe wandering target with offset from obstacles
func GetSafeWanderingTarget(currentX, currentZ, maxRadius, offsetDistance float64, q worldQuery.WorldQuery, rng interface {
	Float64() float64
	Intn(n int) int
}) (float64, float64, bool) {
	const maxAttempts = 10

	for attempt := 0; attempt < maxAttempts; attempt++ {
		angle := rng.Float64() * 2 * math.Pi
		radius := math.Sqrt(rng.Float64()) * maxRadius

		dx := math.Cos(angle) * radius
		dz := math.Sin(angle) * radius

		targetX := currentX + dx
		targetZ := currentZ + dz

		// Check boundaries
		worldWidth, worldHeight := q.GetBoundaries()
		if targetX < 0 || targetX > worldWidth || targetZ < 0 || targetZ > worldHeight {
			continue
		}

		// Find offset position
		if !q.IsPointBlocked(targetX, targetZ) {
			if offsetX, offsetZ, ok := FindOffsetPosition(targetX, targetZ, offsetDistance, q); ok {
				return offsetX, offsetZ, true
			}
		}
	}

	return currentX, currentZ, false
}
