// Package summary serves the summary template.
package summary

import (
	"github.com/pkg/errors"
	"gitlab.com/glatteis/earthwalker/challenge"
	"gitlab.com/glatteis/earthwalker/player"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
)

var summaryPage = template.Must(template.ParseFiles("templates/main_template.html.tmpl", "templates/summary/summary.html.tmpl"))

type guessedPositionsType struct {
	GuessedPosition []float64
	Nickname        string
}

type rankingType struct {
	Nickname            string
	NumPoints           int
	AccumulatedDistance string
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
	DistanceInfo    []distanceType
}

// ServeSummary serves the summary page.
func ServeSummary(w http.ResponseWriter, r *http.Request) {
	session, err := player.GetSessionFromCookie(r)
	if err == player.ErrPlayerSessionNotFound {
		http.Error(w, "you are not authenticated to guess!", http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "some error occured", http.StatusUnprocessableEntity)
		return
	}

	foundChallenge, err := challenge.GetChallenge(session.GameID)
	if err == challenge.ErrChallengeNotFound {
		http.Error(w, "this challenge does not exist!", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "there was some kind of internal error, sorry!", http.StatusUnprocessableEntity)
		return
	}

	if session.Round() < foundChallenge.Settings.NumRounds {
		http.Error(w, "You have not completed every round yet, so you cannot view the summary.", http.StatusUnprocessableEntity)
		return
	}
	ranking := makeRanking(foundChallenge)
	actualPositionsAsFloats, allGuessedPositions := makeMap(foundChallenge)
	distanceInfo := makeDistanceInfo(session)
	err = summaryPage.Execute(w, summaryServeStruct{
		Rankings:        ranking,
		Guesses:         allGuessedPositions,
		ActualPositions: actualPositionsAsFloats,
		DistanceInfo:    distanceInfo,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, "there was some kind of internal error, sorry!", http.StatusUnprocessableEntity)
		return
	}
}

func makeRanking(foundChallenge challenge.Challenge) []rankingType {
	ranking := make([]rankingType, 0)

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
			}
		}
		allGuessedPositions = append(allGuessedPositions, guessedPositions)
	}

	return actualPositionsAsFloats, allGuessedPositions
}

func makeDistanceInfo(session player.PlayerSession) []distanceType {
	distances := make([]distanceType, 0)
	for i, distance := range session.Distances {
		distances = append(distances, distanceType{
			Round:    i + 1,
			Points:   session.Points[i],
			Distance: strconv.FormatFloat(distance, 'f', 2, 64),
		})
	}
	return distances
}
