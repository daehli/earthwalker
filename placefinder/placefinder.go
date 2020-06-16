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
	"gitlab.com/glatteis/earthwalker/domain"
)

type NewMapHandler struct {
	MapStore domain.MapStore
}

func (handler *NewMapHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	mapData := r.FormValue("mapData")
	if mapData == "" {
		http.Error(w, "No mapData received", http.StatusBadRequest)
		return
	}
	newMap, err := mapFromData(mapData)
	if err != nil {
		log.Printf("Failed to create map from data: %v\n", err)
		http.Error("Failed to create map from data.", http.StatusInternalServerError)
		return
	}
	err = handler.MapStore.Insert(newMap)
	if err != nil {
		log.Printf("Failed to insert new map into store: %v\n", err)
		http.Error("Failed to insert new map into store.", http.StatusInternalServerError)
		return
	}
	// TODO: redirect (to new challenge page for this map?)
}

func mapFromData(mapData string) (domain.Map, error) {
	newMap := domain.Map{
		MapID: "",
		Name: "", // TODO: not yet implemented
		Polygon: nil, // TODO: not yet sent from client
		Area: -1, // TODO: not yet implemented
		GraceDistance: 10, // TODO: option not implemented on client side
	}
	err := json.Unmarshal([]byte(mapData), &newMap)
	if err != nil {
		return newMap, fmt.Errorf("Failed to unmarshal newMap from JSON: %v", err)
	}
	// we want to make sure we don't take the ID from the client request
	newMap.MapID = domain.RandAlpha(10)
	return newMap, nil
}

type NewChallengeHandler struct {
	ChallengeStore domain.ChallengeStore
}

func (handler *NewChallengeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	challengeData := r.FormValue("challengeData")
	if challengeData == "" {
		http.Error(w, "No challengeData received", http.StatusBadRequest)
		return
	}
	newChallenge, err := challengeFromData(challengeData)
	if err != nil {
		log.Printf("Failed to create challenge from data: %v\n", err)
		http.Error("Failed to create challenge from data.", http.StatusInternalServerError)
		return
	}
	err = handler.ChallengeStore.Insert(newChallenge)
	if err != nil {
		log.Printf("Failed to insert new challenge into store: %v\n", err)
		http.Error("Failed to insert new challenge into store.", http.StatusInternalServerError)
		return
	}
	// TODO: redirect (to before_start page for this challenge?)
}

func challengeFromData(challengeData string) (domain.Challenge, error) {
	type jsonPoint struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}
	newChallenge := domain.Challenge{
		ChallengeID: domain.RandAlpha(10),
		Places: make([]domain.ChallengePlace, 0),
	}
	var locations []jsonPoint
	if err := json.Unmarshal([]byte(result), &locations); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	// convert from degrees to radians (ffs) and populate challenge.Places
	for i := range locations {
		challenge.Places = append(newChallenge.Places, domain.ChallengePlace{
			ChallengeID: newChallenge.ChallengeID,
			Location: s2.LatLngFromDegrees(locations[i].Lat, locations[i].Lng)
		})
	}
}

type NewChallengeResultHandler struct {
	ChallengeResultStore domain.ChallengeResultStore
}

func (handler *NewChallengeResultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: create new ChallengeResult from form
	// TODO: validate ChallengeResult
	//       nickname not empty
	// TODO: insert into store
	// TODO: redirect to actual game
}
