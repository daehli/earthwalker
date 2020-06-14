// Package `domain` defines domain types which are used throughout the application.
// It imports nothing.
package domain

// Map is a collection of user provided settings from which a Challenge
// is generated.  Formerly challenge/Settings
type Map struct {
}

// MapStore is implemented by structs which provide access to a database
// containing Maps.
type MapStore interface {
}

// Challenge is a list of coordinates of panos.
type Challenge struct {
}

// ChallengeStore is implemented by structs which provide access to a database
// containing Challenges.
type ChallengeStore interface {
}

// ChallengeResult is a player's Guesses in a challenge.
// A ChallengeResult may still be in progress.
// Similar to former player/Session
type ChallengeResult struct {
}

// ChallengeResultStore is implemented by structs which provide access to a
// database containing ChallengeResults.
type ChallengeResultStore interface {
}

// Guess is a guessed location for one pano in a Challenge.
type Guess struct {
}

// GuessStore is implemented by structs which provide access to a database
// containing Guesses.
type GuessStore interface {
}
