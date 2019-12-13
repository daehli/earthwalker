package dynamicpages

import (
	"gitlab.com/glatteis/earthwalker/challenge"
	"gitlab.com/glatteis/earthwalker/player"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var scorePage = template.Must(template.ParseFiles("templates/main_template.html.tmpl", "templates/score/score.html.tmpl"))

type guessedPositionsType struct {
	GuessedPosition []float64
	Nickname        string
}

type scoreServeStruct struct {
	NumPoints        int
	PointsPercent    int
	DistanceKM       string
	GuessedPositions map[string]guessedPositionsType
	ActualPosition   []float64
	LastScorePage    bool
	YourID           string
}

// ServeScores serves the scores page.
func ServeScores(w http.ResponseWriter, r *http.Request) {
	session, err := player.GetSessionFromCookie(r)
	if err == player.PlayerSessionNotFoundError {
		w.Write([]byte("you are not authenticated to guess!"))
		w.WriteHeader(401)
		return
	} else if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	foundChallenge, err := challenge.GetChallenge(session.GameID)
	if err == challenge.ChallengeNotFoundError {
		w.Write([]byte("this challenge does not exist!"))
		w.WriteHeader(404)
		return
	} else if err != nil {
		log.Println(err)
		w.Write([]byte("there was some kind of internal error, sorry!"))
		w.WriteHeader(500)
		return
	}

	if session.Round() <= 1 {
		w.Write([]byte("You have not completed a round yet, you cannot view scores."))
		w.WriteHeader(500)
		return
	}
	actualPosition := foundChallenge.Places[session.Round()-2]
	actualPositionAsFloats := []float64{actualPosition.Lat.Degrees(), actualPosition.Lng.Degrees()}

	guessedPositions := make(map[string]guessedPositionsType)
	for _, guess := range foundChallenge.Guesses[session.Round()-2] {
		guessedPositions[guess.PlayerID] = guessedPositionsType{
			GuessedPosition: []float64{guess.GuessLocation.Lat.Degrees(), guess.GuessLocation.Lng.Degrees()},
			Nickname:        guess.PlayerNickname,
		}
	}

	toServe := scoreServeStruct{
		NumPoints:        session.Points[session.Round()-2],
		PointsPercent:    session.Points[session.Round()-2] / (5000 / 100),
		DistanceKM:       strconv.FormatFloat(session.Distances[session.Round()-2], 'f', 2, 64),
		GuessedPositions: guessedPositions,
		ActualPosition:   actualPositionAsFloats,
		LastScorePage:    session.Round()-1 == foundChallenge.Settings.NumRounds,
		YourID:           session.UniqueIdentifier,
	}

	err = scorePage.Execute(w, toServe)
	if err != nil {
		log.Println(err)
		w.Write([]byte("there was some kind of internal error, sorry!"))
		w.WriteHeader(500)
		return
	}
}
