// Package scores calculates scores.
package scores

import (
	"math"

	"github.com/golang/geo/s2"
)

const earthRadius = 6371

const maxScore = 5000
const graceDistance = 0.025 // perfect scores will be given within this distance (km)

// effectively, score will be divided by decayBase every decayDistance km
const decayBase = 2
const decayDistance = 1070

// CalculateScoreAndDistance calculates the score and distance for a guessed location (and its actual location).
func CalculateScoreAndDistance(actualLocation s2.LatLng, guessLocation s2.LatLng) (int, float64) {
	distance := actualLocation.Distance(guessLocation).Radians() * earthRadius
	if distance < graceDistance {
		return maxScore, distance
	}
	factor := math.Pow(decayBase, -float64(distance)/decayDistance)
	points := int(factor * maxScore)
	return points, distance
}
