// Package summary serves the summary template.
package summary

import (
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/pkg/errors"
	"gitlab.com/glatteis/earthwalker/challenge"
	"gitlab.com/glatteis/earthwalker/player"
	"gitlab.com/glatteis/earthwalker/util"
)

var summaryPage = template.Must(template.ParseFiles(util.StaticPath()+"/templates/main_template.html.tmpl", util.StaticPath()+"/templates/summary/summary.html.tmpl"))

type guessedPositionsType struct {
	GuessedPosition []float64
	Nickname        string
	Color           int
}

type rankingType struct {
	Nickname            string
	NumPoints           int
	AccumulatedDistance string
	Color               int
}

type distanceType struct {
	Round    int
	Points   int
	Distance string
}

type summaryServeStruct struct {
	Rankings        []rankingType
	Guesses         []map[string]guessedPositionsType
	ActualPositions [][]float64
	HasDistanceInfo bool
	GameID          string
	DistanceInfo    []distanceType
}

// ServeSummary serves the summary page.
func ServeSummary(w http.ResponseWriter, r *http.Request) {
	// Look if there is a challenge key and password present.
	// If yes, see if there is a session present, if yes, serve the summary. If no session is present,
	// serve it without the distance info.
	// If no, see if there is a session present, if yes, redirect to page with password and game id.
	// If none are present, reject the user.
	session, sessionErr := player.GetSessionFromCookie(r)
	var foundChallenge challenge.Challenge
	hasSession := sessionErr == nil
	if hasSession {
		var err error
		foundChallenge, err = challenge.GetChallenge(session.GameID)
		if err == challenge.ErrChallengeNotFound {
			http.Error(w, "this challenge does not exist!", http.StatusNotFound)
			return
		} else if err != nil {
			log.Println(err)
			http.Error(w, "there was some kind of internal error, sorry!", http.StatusUnprocessableEntity)
			return
		}
	}

	challengeKeys, ok := r.URL.Query()["c"]
	if !ok || len(challengeKeys) == 0 || challengeKeys[0] == "" {
		if !hasSession {
			http.Error(w, "you cannot view a summary without having played the game or having an id", http.StatusUnauthorized)
			return
		}
	}

	enteredPasswords, ok := r.URL.Query()["p"]
	if !ok || len(enteredPasswords) == 0 || enteredPasswords[0] == "" {
		if !hasSession {
			http.Error(w, "you cannot view a summary without having played the game or having a password", http.StatusUnauthorized)
			return
		}
		// At this point we know there's someone with a valid session but no game id / password yet calling.

		// Reject the user if they have a session but have not completed the game yet.
		if session.Round() <= foundChallenge.Settings.NumRounds {
			http.Error(w, "you have not completed every round yet, so you cannot view the summary", http.StatusUnprocessableEntity)
			return
		}

		// Redirect them to the url with id and password.
		http.Redirect(w, r, "/summary?c="+foundChallenge.UniqueIdentifier+"&p="+foundChallenge.SummaryPassword, http.StatusFound)
		return
	}

	// Treat having a session but having not completed the game like having no session
	hasSession = hasSession && !(session.Round() <= foundChallenge.Settings.NumRounds) && (session.GameID == foundChallenge.UniqueIdentifier)

	// Load the challenge if they don't have a session.
	if !hasSession {
		var err error
		foundChallenge, err = challenge.GetChallenge(challengeKeys[0])
		if err != nil {
			http.Error(w, "that challenge does not exist", http.StatusUnauthorized)
		}
	}

	// Check the password
	if enteredPasswords[0] != foundChallenge.SummaryPassword {
		http.Error(w, "incorrect password!", http.StatusUnauthorized)
		return
	}

	// At this point, we know the person is authorized to view the summary.

	ranking := makeRanking(foundChallenge)
	actualPositionsAsFloats, allGuessedPositions := makeMap(foundChallenge)
	var distanceInfo []distanceType
	if hasSession {
		distanceInfo = makeDistanceInfo(session)
	}
	err := summaryPage.Execute(w, summaryServeStruct{
		Rankings:        ranking,
		Guesses:         allGuessedPositions,
		ActualPositions: actualPositionsAsFloats,
		HasDistanceInfo: hasSession,
		DistanceInfo:    distanceInfo,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, "there was some kind of internal error, sorry!", http.StatusUnprocessableEntity)
		return
	}
}

func makeRanking(foundChallenge challenge.Challenge) []rankingType {
	var ranking []rankingType

	for _, playerThatCompleted := range foundChallenge.Guesses[len(foundChallenge.Guesses)-1] {
		completedSession, err := player.LoadPlayerSession(playerThatCompleted.PlayerID)
		if err != nil {
			log.Println(errors.Wrap(err, "while loading a player that should have guessed"))
		}
		var sumPoints int
		var sumDistance float64
		for _, p := range completedSession.Points {
			sumPoints += p
		}
		for _, distance := range completedSession.Distances {
			sumDistance += distance
		}
		ranking = append(ranking, rankingType{
			Nickname:            completedSession.Nickname,
			NumPoints:           sumPoints,
			AccumulatedDistance: strconv.FormatFloat(sumDistance, 'f', 2, 64),
			Color:               completedSession.IconColor,
		})
	}

	sort.Slice(ranking, func(i, j int) bool {
		return ranking[i].NumPoints >= ranking[j].NumPoints
	})
	return ranking
}

func makeMap(foundChallenge challenge.Challenge) ([][]float64, []map[string]guessedPositionsType) {
	actualPositions := foundChallenge.Places
	actualPositionsAsFloats := [][]float64{}
	for _, position := range actualPositions {
		actualPositionsAsFloats = append(actualPositionsAsFloats, []float64{position.Lat.Degrees(), position.Lng.Degrees()})
	}

	allGuessedPositions := []map[string]guessedPositionsType{}

	for round := range foundChallenge.Guesses {
		guessedPositions := make(map[string]guessedPositionsType)
		for _, guess := range foundChallenge.Guesses[round] {
			guessedPositions[guess.PlayerID] = guessedPositionsType{
				GuessedPosition: []float64{guess.GuessLocation.Lat.Degrees(), guess.GuessLocation.Lng.Degrees()},
				Nickname:        guess.PlayerNickname,
				Color:           guess.PlayerColor,
			}
		}
		allGuessedPositions = append(allGuessedPositions, guessedPositions)
	}

	return actualPositionsAsFloats, allGuessedPositions
}

func makeDistanceInfo(session player.Session) []distanceType {
	var distances []distanceType
	for i, distance := range session.Distances {
		distances = append(distances, distanceType{
			Round:    i + 1,
			Points:   session.Points[i],
			Distance: strconv.FormatFloat(distance, 'f', 2, 64),
		})
	}
	return distances
}
