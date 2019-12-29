// Package modifyfrontend serves the modify_frontend.js template. Yes, it imports text/template on purpose
package modifyfrontend

import (
	"gitlab.com/glatteis/earthwalker/challenge"
	"gitlab.com/glatteis/earthwalker/player"
	"log"
	"net/http"
	"text/template"
	"time"
)

type modifyServeStruct struct {
	TimerEnabled   bool
	TimerDuration  int
	LabeledMinimap bool
}

var modifyScript = template.Must(template.ParseFiles("templates/modify_frontend/modify.js.tmpl"))

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
	}

	toServe.LabeledMinimap = game.Settings.LabeledMinimap

	err = modifyScript.Execute(w, toServe)
	if err != nil {
		log.Println(err)
		http.Error(w, "there was some kind of internal error, sorry!", http.StatusUnprocessableEntity)
		return
	}
}
