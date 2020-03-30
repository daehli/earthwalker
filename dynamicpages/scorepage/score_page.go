// Package scorepage serves the score_page template.
package scorepage

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"gitlab.com/glatteis/earthwalker/challenge"
	"gitlab.com/glatteis/earthwalker/player"
	"gitlab.com/glatteis/earthwalker/util"
)

var scorePage = template.Must(template.ParseFiles(util.AppPath()+"/templates/main_template.html.tmpl", util.AppPath()+"/templates/score/score.html.tmpl"))

type guessedPositionsType struct {
	GuessedPosition []float64
	Nickname        string
	Color           int
}

type scoreServeStruct struct {
	NumPoints        int
	PointsPercent    int
	DistanceKM       string
	GuessedPositions map[string]guessedPositionsType
	ActualPosition   []float64
	LastScorePage    bool
	YourColor        int
}

// ServeScores serves the scores page.
func ServeScores(w http.ResponseWriter, r *http.Request) {
	session, err := player.GetSessionFromCookie(r)
	if err == player.ErrPlayerSessionNotFound {
		http.Error(w, "you are not authenticated to guess!", http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "some error happened", http.StatusUnprocessableEntity)
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

	if session.Round() <= 1 {
		http.Error(w, "You have not completed a round yet, you cannot view scores.", http.StatusUnprocessableEntity)
		return
	}
	actualPosition := foundChallenge.Places[session.Round()-2]
	actualPositionAsFloats := []float64{actualPosition.Lat.Degrees(), actualPosition.Lng.Degrees()}

	guessedPositions := make(map[string]guessedPositionsType)
	for _, guess := range foundChallenge.Guesses[session.Round()-2] {
		guessedPositions[guess.PlayerID] = guessedPositionsType{
			GuessedPosition: []float64{guess.GuessLocation.Lat.Degrees(), guess.GuessLocation.Lng.Degrees()},
			Nickname:        guess.PlayerNickname,
			Color:           guess.PlayerColor,
		}
	}

	toServe := scoreServeStruct{
		NumPoints:        session.Points[session.Round()-2],
		PointsPercent:    session.Points[session.Round()-2] / (5000 / 100),
		DistanceKM:       strconv.FormatFloat(session.Distances[session.Round()-2], 'f', 2, 64),
		GuessedPositions: guessedPositions,
		ActualPosition:   actualPositionAsFloats,
		LastScorePage:    session.Round()-1 == foundChallenge.Settings.NumRounds,
		YourColor:        session.IconColor,
	}

	w.Header().Set("Cache-Control", "no-cache")
	err = scorePage.Execute(w, toServe)
	if err != nil {
		log.Println(err)
		http.Error(w, "there was some kind of internal error, sorry!", http.StatusUnprocessableEntity)
		return
	}
}
