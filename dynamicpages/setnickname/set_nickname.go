// Package setnickname serves the set_nickname template.
package setnickname

import (
	"html/template"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"gitlab.com/glatteis/earthwalker/player"
	"gitlab.com/glatteis/earthwalker/util"
)

var setNickname = template.Must(template.ParseFiles(util.StaticPath()+"/templates/main_template.html.tmpl", util.StaticPath()+"/templates/set_nickname/set_nickname.html.tmpl"))

type nicknameServeStruct struct {
	GameID string
}

// ServeSetNickname serves the set_nickname template.
func ServeSetNickname(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		setNicknamePost(w, r)
		return
	}

	challengeKey, ok := r.URL.Query()["c"]
	// This is probably what they call "user error"
	if !ok || len(challengeKey) == 0 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	actualKey := challengeKey[0]
	err := setNickname.Execute(w, nicknameServeStruct{
		GameID: actualKey,
	})
	if err != nil {
		log.Println(errors.Wrap(err, "while executing the set_nickname template"))
		http.Error(w, "Could not serve you a template for some reason, sorry!", http.StatusUnprocessableEntity)
		return
	}
}

func setNicknamePost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	nickname := r.FormValue("nickname")
	if nickname == "" {
		http.Error(w, "Nickname cannot be empty!", http.StatusUnprocessableEntity)
		return
	}
	gameID := r.FormValue("game_id")

	player.WriteNicknameAndSession(w, r, nickname)

	http.Redirect(w, r, "/game?c="+gameID, http.StatusFound)
}
