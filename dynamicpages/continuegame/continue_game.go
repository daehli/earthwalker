// Package continueGame serves the template for asking whether a new game should be started.
package continuegame

import (
	"html/template"
	"log"
	"net/http"
)

var continueGame = template.Must(template.ParseFiles("templates/main_template.html.tmpl", "templates/continue_game/continue_game.html.tmpl"))

// ServeContinueGame serves the continuegame template.
func ServeContinueGame(w http.ResponseWriter, r *http.Request) {
	err := continueGame.Execute(w, struct{}{})
	if err != nil {
		log.Println(err)
		http.Error(w, "there was some kind of internal error, sorry!", http.StatusUnprocessableEntity)
		return
	}
}
