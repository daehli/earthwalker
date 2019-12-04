package challenge

import (
	"encoding/json"
	"github.com/golang/geo/s2"
	"gitlab.com/glatteis/earthwalker/player"
	"io/ioutil"
	"log"
	"math"
	"net/http"
)

const earthRadius = 6371

// Guess serves the post request that is sent when one guesses.
func Guess(w http.ResponseWriter, r *http.Request) {
	session, err := player.GetSessionFromCookie(r)
	if err == player.PlayerSessionNotFoundError {
		w.Write([]byte("you are not authenticated to guess!"))
		w.WriteHeader(401)
		return
	} else if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	foundChallenge, err := GetChallenge(session.GameID)
	if err == ChallengeNotFoundError {
		w.Write([]byte("this challenge does not exist!"))
		w.WriteHeader(404)
		return
	} else if err != nil {
		log.Println(err)
		w.Write([]byte("there was some kind of internal error, sorry!"))
		w.WriteHeader(500)
		return
	}

	actualLocation := foundChallenge.Places[session.Round()-1]

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.Write([]byte("there was some kind of internal error, sorry!"))
		w.WriteHeader(500)
		return
	}

	guessLocation, err := parseGuessLocation(body)
	if err != nil {
		log.Println(err)
		w.Write([]byte("there was some kind of error while parsing the input json."))
		w.WriteHeader(500)
		return
	}

	distance := actualLocation.Distance(guessLocation).Radians() * earthRadius
	maxDistance := earthRadius * math.Pi
	points := int(5000 - ((float64(distance) / maxDistance) * 5000))

	session.Points = append(session.Points, points)
	session.GuessedPositions = append(session.GuessedPositions, []float64{guessLocation.Lat.Degrees(), guessLocation.Lng.Degrees()})
	session.Distances = append(session.Distances, distance)

	log.Println("Guessed!")
	log.Println(session.Points)

	err = player.StorePlayerSession(session)
	if err != nil {
		log.Println(err)
		w.Write([]byte("there was some kind of internal error, sorry!"))
		w.WriteHeader(500)
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
