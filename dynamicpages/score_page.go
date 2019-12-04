package dynamicpages

import (
	"gitlab.com/glatteis/earthwalker/challenge"
	"gitlab.com/glatteis/earthwalker/player"
	"html/template"
	"log"
	"net/http"
)

var scorePage = template.Must(template.ParseFiles("templates/main_template.html.tmpl", "templates/score/score.html.tmpl"))

type serveStruct struct {
	NumPoints       int
	DistanceKM      float64
	GuessedPosition []float64
	ActualPosition  []float64
}

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

	log.Println(session.Round())

	actualPosition := foundChallenge.Places[session.Round()-2]
	actualPositionAsFloats := []float64{actualPosition.Lat.Degrees(), actualPosition.Lng.Degrees()}

	toServe := serveStruct{
		NumPoints:       session.Points[session.Round()-2],
		DistanceKM:      session.Distances[session.Round()-2],
		GuessedPosition: session.GuessedPositions[session.Round()-2],
		ActualPosition:  actualPositionAsFloats,
	}

	err = scorePage.Execute(w, toServe)
	if err != nil {
		log.Println(err)
		w.Write([]byte("there was some kind of internal error, sorry!"))
		w.WriteHeader(500)
		return
	}
}
