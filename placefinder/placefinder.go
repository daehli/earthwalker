package placefinder

import (
	"encoding/json"
	"errors"
	"github.com/golang/geo/s2"
	"gitlab.com/glatteis/earthwalker/challenge"
	"log"
	"net/http"
	"strconv"
)

func RespondToPoints(w http.ResponseWriter, r *http.Request) {
	type jsonPoint struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}

	r.ParseForm()
	result := r.FormValue("result")

	var content []jsonPoint

	if err := json.Unmarshal([]byte(result), &content); err != nil {
		http.Error(w, err.Error(), 403)
		return
	}

	locations := make([]s2.LatLng, len(content))
	for i := range content {
		locations[i] = s2.LatLngFromDegrees(content[i].Lat, content[i].Lng)
	}

	nickname := r.FormValue("nickname")
	if nickname == "" {
		w.Write([]byte("Nickname cannot be empty!"))
		w.WriteHeader(422)
		return
	}

	challenge.WriteNicknameAndSession(w, r, nickname)

	settings, err := createSettingsFromForm(r)

	if err != nil {
		w.Write([]byte("There was something wrong with your parameters: " + err.Error()))
		return
	}

	resultingChallenge, err := challenge.NewChallenge(locations, settings)

	if err != nil {
		log.Println(err)
		w.Write([]byte("Internal server error! Please contact an administrator or something"))
		return
	}

	http.Redirect(w, r, "/game?c="+resultingChallenge.UniqueIdentifier, 302)
}

func createSettingsFromForm(r *http.Request) (challenge.ChallengeSettings, error) {
	var settings challenge.ChallengeSettings
	r.ParseForm()

	numRoundsStr := r.FormValue("rounds")
	roundsAsInt, err := strconv.Atoi(numRoundsStr)
	if err != nil {
		return challenge.ChallengeSettings{}, errors.New("rounds is not an integer!")
	}
	settings.NumRounds = roundsAsInt

	return settings, nil
}
