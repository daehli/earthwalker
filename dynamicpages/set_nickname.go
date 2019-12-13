package dynamicpages

import (
	"github.com/pkg/errors"
	"gitlab.com/glatteis/earthwalker/challenge"
	"html/template"
	"log"
	"net/http"
)

var setNickname = template.Must(template.ParseFiles("templates/main_template.html.tmpl", "templates/set_nickname/set_nickname.html.tmpl"))

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
		http.Redirect(w, r, "/", 302)
		return
	}
	actualKey := challengeKey[0]
	err := setNickname.Execute(w, nicknameServeStruct{
		GameID: actualKey,
	})
	if err != nil {
		log.Println(errors.Wrap(err, "while executing the set_nickname template"))
		w.Write([]byte("Could not serve you a template for some reason, sorry!"))
		w.WriteHeader(500)
		return
	}
}

func setNicknamePost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	nickname := r.FormValue("nickname")
	if nickname == "" {
		w.Write([]byte("Nickname cannot be empty!"))
		w.WriteHeader(422)
		return
	}
	gameID := r.FormValue("game_id")

	challenge.WriteNicknameAndSession(w, r, nickname)

	http.Redirect(w, r, "/game?c="+gameID, 302)
}
