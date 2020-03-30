// Package beforestart serves the before_start template.
package beforestart

import (
	"html/template"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"gitlab.com/glatteis/earthwalker/util"
)

var beforeStart = template.Must(template.ParseFiles(util.AppPath()+"/templates/main_template.html.tmpl", util.AppPath()+"/templates/before_start/before_start.html.tmpl"))

type beforeStartServeStruct struct {
	GameID string
}

// ServeBeforeStart serves the before_start template.
func ServeBeforeStart(w http.ResponseWriter, r *http.Request) {
	challengeKey, ok := r.URL.Query()["c"]
	// This is probably what they call "user error"
	if !ok || len(challengeKey) == 0 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	actualKey := challengeKey[0]
	err := beforeStart.Execute(w, beforeStartServeStruct{
		GameID: actualKey,
	})
	if err != nil {
		log.Println(errors.Wrap(err, "while executing the before_start template"))
		http.Error(w, "Could not serve you a template for some reason, sorry!", http.StatusUnprocessableEntity)
		return
	}
}
