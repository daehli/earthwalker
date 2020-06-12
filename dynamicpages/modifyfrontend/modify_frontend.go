// Package modifyfrontend serves the modify_frontend.js template. Yes, it imports text/template on purpose
package modifyfrontend

import (
	"log"
	"net/http"
	"text/template"
	"time"

	"gitlab.com/glatteis/earthwalker/challenge"
	"gitlab.com/glatteis/earthwalker/player"
	"gitlab.com/glatteis/earthwalker/util"
)

type modifyServeStruct struct {
	TimerEnabled   bool
	TimerDuration  int
	LabeledMinimap bool
	RoundNumber    int
	NumberOfRounds int
	GameID         string
}

var modifyScript = template.Must(template.ParseFiles(util.StaticPath() + "/templates/modify_frontend/modify.js.tmpl"))

// ServeModifyFrontend serves the modify_frontend.js template.
func ServeModifyFrontend(w http.ResponseWriter, r *http.Request) {
	session, err := player.GetSessionFromCookie(r)
	if err != nil {
		http.Error(w, "not authorized", http.StatusUnauthorized)
	}

	game, err := challenge.GetChallenge(session.GameID)
	if err != nil {
		http.Error(w, "the game does not exist", http.StatusNotFound)
	}

	var toServe modifyServeStruct

	if game.Settings.TimerDuration != nil {
		toServe.TimerEnabled = true
		duration := *game.Settings.TimerDuration
		timeStarted := session.TimeStarted
		if timeStarted == nil {
			log.Println("Error: game has timer, but player has not stared game!")
			http.Error(w, "internal server error!", http.StatusInternalServerError)
		}
		alreadyPassed := time.Since(*timeStarted)
		toServe.TimerDuration = int((duration - alreadyPassed) / time.Second)
		if toServe.TimerDuration < 1 {
			toServe.TimerDuration = 1
		}
	}

	toServe.LabeledMinimap = game.Settings.LabeledMinimap
	toServe.RoundNumber = session.Round()
	toServe.NumberOfRounds = game.Settings.NumRounds
	toServe.GameID = game.UniqueIdentifier

	err = modifyScript.Execute(w, toServe)
	if err != nil {
		log.Println(err)
		http.Error(w, "there was some kind of internal error, sorry!", http.StatusUnprocessableEntity)
		return
	}
}
