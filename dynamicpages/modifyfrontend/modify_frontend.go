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
	TimerEnabled  bool
	TimerDuration int
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
		http.Error(w, "the game does not exist", http.StatusInternalServerError)
	}

	var toServe modifyServeStruct

	if game.Settings.TimerDuration != nil {
		toServe.TimerEnabled = true
		toServe.TimerDuration = int(*game.Settings.TimerDuration / time.Second)
	}

	err = modifyScript.Execute(w, toServe)
	if err != nil {
		log.Println(err)
		http.Error(w, "there was some kind of internal error, sorry!", http.StatusUnprocessableEntity)
		return
	}
}
