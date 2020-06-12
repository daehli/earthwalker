// Package placefinder serves the page that responds to found places.
package placefinder

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang/geo/s2"
	"gitlab.com/glatteis/earthwalker/challenge"
	"gitlab.com/glatteis/earthwalker/player"
)

// RespondToPoints responds to found places.
func RespondToPoints(w http.ResponseWriter, r *http.Request) {
	type jsonPoint struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}

	r.ParseForm()
	result := r.FormValue("result")

	var content []jsonPoint

	if err := json.Unmarshal([]byte(result), &content); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	locations := make([]s2.LatLng, len(content))
	for i := range content {
		locations[i] = s2.LatLngFromDegrees(content[i].Lat, content[i].Lng)
	}

	nickname := r.FormValue("nickname")
	if nickname == "" {
		http.Error(w, "Nickname cannot be empty!", http.StatusUnprocessableEntity)
		return
	}

	player.WriteNicknameAndSession(w, r, nickname)

	settings, err := createSettingsFromForm(r)

	if err != nil {
		http.Error(w, "There was something wrong with your parameters: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	resultingChallenge, err := challenge.NewChallenge(locations, settings)

	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error! Please contact an administrator or something", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/beforestart?c="+resultingChallenge.UniqueIdentifier, http.StatusFound)
}

func createSettingsFromForm(r *http.Request) (challenge.Settings, error) {
	var settings challenge.Settings
	r.ParseForm()

	showLabelsStr := r.FormValue("show-labels")
	if showLabelsStr != "" {
		settings.LabeledMinimap = true
	}

	numRoundsStr := r.FormValue("rounds")
	roundsAsInt, err := strconv.Atoi(numRoundsStr)
	if err != nil {
		return challenge.Settings{}, errors.New("rounds is not an integer")
	}
	if roundsAsInt == 0 {
		return challenge.Settings{}, errors.New("rounds must not be zero")
	}
	settings.NumRounds = roundsAsInt

	var incorrectFormat bool
	roundDurationStr := r.FormValue("time")
	if roundDurationStr != "" {
		twoNumbers := strings.Split(roundDurationStr, ":")
		if len(twoNumbers) != 2 {
			incorrectFormat = true
			goto done
		}
		minutes, err := strconv.Atoi(twoNumbers[0])
		if err != nil {
			incorrectFormat = true
			goto done
		}
		seconds, err := strconv.Atoi(twoNumbers[1])
		if err != nil {
			incorrectFormat = true
			goto done
		}
		duration := time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second
		if duration <= 0 {
			incorrectFormat = true
			goto done
		}
		settings.TimerDuration = &duration
	}
done:
	if incorrectFormat {
		return challenge.Settings{}, errors.New("time is in an incorrect format")
	}

	return settings, nil
}
