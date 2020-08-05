package api

// TODO: improve logging (we probably want to output actual structs/IDs
// 	     when something goes wrong)

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"

	"gitlab.com/glatteis/earthwalker/domain"
)

type Root struct {
	Config               domain.Config
	MapStore             domain.MapStore
	ChallengeStore       domain.ChallengeStore
	ChallengeResultStore domain.ChallengeResultStore

	ConfigHandler     Config
	MapsHandler       Maps
	ChallengesHandler Challenges
	ResultsHandler    Results
	GuessesHandler    Guesses
}

func (handler Root) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("handling api request...")
	log.Println(r.URL.Path)
	head, tail := shiftPath(r.URL.Path)
	r.URL.Path = tail
	switch head {
	case "config":
		handler.ConfigHandler.ServeHTTP(w, r)
	case "maps":
		handler.MapsHandler.ServeHTTP(w, r)
	case "challenges":
		handler.ChallengesHandler.ServeHTTP(w, r)
	case "results":
		handler.ResultsHandler.ServeHTTP(w, r)
	case "guesses":
		handler.GuessesHandler.ServeHTTP(w, r)
	default:
		sendError(w, fmt.Sprintf("API endpoint '%s' does not exist.", head), http.StatusNotFound)
		return
	}
}

// sendError text as JSON
func sendError(w http.ResponseWriter, text string, status int) {
	respJSON := "{error: \"" + text + "\"}"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err := w.Write([]byte(respJSON))
	if err != nil {
		log.Printf("Error writing response: %v\n", err)
	}
}

func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
