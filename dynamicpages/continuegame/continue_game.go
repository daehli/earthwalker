// Package continueGame serves the template for asking whether a new game should be started.
package continuegame

import (
	"gitlab.com/glatteis/earthwalker/player"
	"html/template"
	"log"
	"net/http"
)

type modifyContinueGameStruct struct {
	Nickname string
}

var continueGame = template.Must(template.ParseFiles("templates/main_template.html.tmpl", "templates/continue_game/continue_game.html.tmpl"))

// ServeContinueGame serves the continuegame template.
func ServeContinueGame(w http.ResponseWriter, r *http.Request, session player.PlayerSession) {
	var toContinueGame modifyContinueGameStruct
	toContinueGame.Nickname = session.Nickname
	err := continueGame.Execute(w, toContinueGame)
	if err != nil {
		log.Println(err)
		http.Error(w, "there was some kind of internal error, sorry!", http.StatusUnprocessableEntity)
		return
	}
}
