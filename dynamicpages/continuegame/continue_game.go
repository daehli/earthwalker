// Package continuegame serves the template for asking whether a new game should be started.
package continuegame

import (
	"html/template"
	"log"
	"net/http"

	"gitlab.com/glatteis/earthwalker/config"
)

type modifyContinueGameStruct struct {
	Nickname string
}

var continueGame = template.Must(template.ParseFiles(config.Env.EarthwalkerStaticPath+"/templates/main_template.html.tmpl", config.Env.EarthwalkerStaticPath+"/templates/continue_game/continue_game.html.tmpl"))

// ServeContinueGame serves the continue_game template. It takes the nickname of the current player.
func ServeContinueGame(w http.ResponseWriter, r *http.Request, nickname string) {
	toContinueGame := modifyContinueGameStruct{
		Nickname: nickname,
	}
	err := continueGame.Execute(w, toContinueGame)
	if err != nil {
		log.Println(err)
		http.Error(w, "there was some kind of internal error, sorry!", http.StatusUnprocessableEntity)
		return
	}
}
