package challenge

import (
	"encoding/json"
	"github.com/golang/geo/s2"
	"github.com/pkg/errors"
	"gitlab.com/glatteis/earthwalker/player"
	"gitlab.com/glatteis/earthwalker/scores"
	"io/ioutil"
	"log"

	"net/http"
)

// HandleGuess serves the post request that is sent when one guesses.
func HandleGuess(w http.ResponseWriter, r *http.Request) {
	session, err := player.GetSessionFromCookie(r)
	if err == player.ErrPlayerSessionNotFound {
		http.Error(w, "you are not authenticated to guess!", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "there was some kind of internal error, sorry!", http.StatusInternalServerError)
		return
	}

	foundChallenge, err := GetChallenge(session.GameID)
	if err == ErrChallengeNotFound {
		http.Error(w, "this challenge does not exist!", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "there was some kind of internal error, sorry!", http.StatusUnprocessableEntity)
		return
	}

	actualLocation := foundChallenge.Places[session.Round()-1]

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "there was some kind of internal error, sorry!", http.StatusUnprocessableEntity)
		return
	}

	guessLocation, err := parseGuessLocation(body)
	if err != nil {
		log.Println(err)
		http.Error(w, "there was some kind of error while parsing the input json.", http.StatusUnprocessableEntity)
		return
	}

	points, distance := scores.CalculateScoreAndDistance(actualLocation, guessLocation)

	// map guess longitudes so that the distance appearing on the map seems shortest
	if guessLocation.Lng.Degrees() - 180 > actualLocation.Lng.Degrees() {
		guessLocation = s2.LatLngFromDegrees(guessLocation.Lat.Degrees(), guessLocation.Lng.Degrees() - 360)
	} else if guessLocation.Lng.Degrees() + 180 < actualLocation.Lng.Degrees() {
		guessLocation = s2.LatLngFromDegrees(guessLocation.Lat.Degrees(), guessLocation.Lng.Degrees() + 360)
	}

	foundChallenge.Guesses[session.Round()-1] = append(foundChallenge.Guesses[session.Round()-1], Guess{
		GuessLocation:  guessLocation,
		PlayerID:       session.UniqueIdentifier,
		PlayerNickname: session.Nickname,
	})

	session.TimeStarted = nil
	session.Points = append(session.Points, points)
	session.GuessedPositions = append(session.GuessedPositions, []float64{guessLocation.Lat.Degrees(), guessLocation.Lng.Degrees()})
	session.Distances = append(session.Distances, distance)

	err = player.StorePlayerSession(session)
	if err != nil {
		log.Println(errors.Wrap(err, "while storing the session"))
		http.Error(w, "There was an error while storing your session.", http.StatusUnprocessableEntity)
		return
	}

	err = StoreChallenge(foundChallenge)
	if err != nil {
		log.Println(errors.Wrap(err, "while storing the challenge"))
		http.Error(w, "There was an error while updating the challenge.", http.StatusUnprocessableEntity)
		return
	}
}

func parseGuessLocation(body []byte) (s2.LatLng, error) {
	type guessLocation struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lng"`
	}

	var target guessLocation

	if err := json.Unmarshal(body, &target); err != nil {
		return s2.LatLng{}, err
	}

	return s2.LatLngFromDegrees(target.Lat, target.Lon), nil
}
