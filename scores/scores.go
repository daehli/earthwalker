// Package scores calculates scores.
package scores

import (
	"github.com/golang/geo/s2"
	"math"
)

const earthRadius = 6371

func CalculateScoreAndDistance(actualLocation s2.LatLng, guessLocation s2.LatLng) (int, float64) {
	distance := actualLocation.Distance(guessLocation).Radians() * earthRadius
	maxDistance := earthRadius * math.Pi
	factor := math.Pow(1-(float64(distance)/maxDistance), 5)
	points := int(factor * 5000)
	return points, distance
}
