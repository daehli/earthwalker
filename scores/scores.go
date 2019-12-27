// Package scores calculates scores.
package scores

import (
	"github.com/golang/geo/s2"
	"math"
)

const earthRadius = 6371

// CalculateScoreAndDistance calculates the score and distance for a guessed location (and its actual location).
func CalculateScoreAndDistance(actualLocation s2.LatLng, guessLocation s2.LatLng) (int, float64) {
	distance := actualLocation.Distance(guessLocation).Radians() * earthRadius
	factor := math.Pow(2, -float64(distance)/1070)
	points := int(factor * 5001)
	return points, distance
}
