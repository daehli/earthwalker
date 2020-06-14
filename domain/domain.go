// Package `domain` defines domain types which are used throughout the application.
// It imports only external packages which compose domain types.
// TODO: Might be nice to rename this package "earthwalker", but conflicts with
//       the default executable filename.
package domain

import (
	"time"

	"github.com/golang/geo/s2"
)

// == Domain Enums ========

// PanoConnectedness is the enum representing that Map option
type PanoConnectedness int

const (
	ConnectedAny    PanoConnectedness = iota
	ConnectedAlways                   // include only panos with at least one adjacent pano
	ConnectedNever                    // include only panos with no adjacent panos
)

func (c PanoConnectedness) String() string {
	return [...]string{"any", "always", "never"}[c]
}

// PanoCopyright is the enum representing that Map option
type PanoCopyright int

const (
	CopyrightAny        PanoCopyright = iota
	CopyrightGoogle                   // include only panos with Google in the copyright
	CopyrightThirdParty               // include only panos without Google in the copyright
)

func (c PanoCopyright) String() string {
	return [...]string{"any", "Google only", "third party only"}[c]
}

// PanoSource is the enum representing that Map option
type PanoSource int

const (
	SourceAny      PanoSource = iota
	SourceOutdoors            // query the streetview API only for outdoors panos
)

// == Domain Types and Stores ========

// Map is a collection of user provided settings from which a Challenge
// is generated.  Formerly challenge/Settings
type Map struct {
	MapID         string
	Polygon       string  // geoJSON string bounding the game area(s)
	Area          float32 // area in Polygon
	NumRounds     int
	TimeLimit     *time.Duration // time limit per round
	MinDensity    int            // minimum population density
	MaxDensity    int            // maximum population density
	Connectedness PanoConnectedness
	Copyright     PanoCopyright
	Source        PanoSource
	ShowLabels    bool // whether to display place labels on the in-game minimap
}

// MapStore is implemented by structs which provide access to a database
// containing Maps.
type MapStore interface {
}

// Challenge is a list of coordinates of panos.
type Challenge struct {
	ChallengeID string
	MapID       string
	Places      []ChallengePlace `db:"-"`
}

// ChallengeStore is implemented by structs which provide access to a database
// containing Challenges.
type ChallengeStore interface {
}

// ChallengePlace is the location of a pano.
// (May contain FOV, heading, etc. in the future.)
type ChallengePlace struct {
	ChallengeID string
	Location    s2.LatLng
}

// ChallengePlaceStore is implemented by structs which provide access to a
// database containing ChallengePlaces.
type ChallengePlaceStore struct {
}

// ChallengeResult is a player's Guesses in a challenge.
// A ChallengeResult may still be in progress.
// Similar to former player/Session
type ChallengeResult struct {
	ChallengeResultID string
	ChallengeID       string

	// these could potentially be replaced with UserID if we were to
	// implement user auth/accounts
	Nickname  string
	Icon      int
	StartTime *time.Time

	Guesses []Guess `db:"-"`
}

// ChallengeResultStore is implemented by structs which provide access to a
// database containing ChallengeResults.
type ChallengeResultStore interface {
}

// Guess is a guessed location for one pano in a Challenge.
type Guess struct {
	ChallengeResultID string
	Location          s2.LatLng
}

// GuessStore is implemented by structs which provide access to a database
// containing Guesses.
type GuessStore interface {
}
