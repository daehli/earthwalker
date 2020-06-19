// handlers in this file create and store new structs before the game begins
// (Map, Challenge, ChallengeResult, etc.)
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/golang/geo/s2"
	"gitlab.com/glatteis/earthwalker/domain"
)

type NewMap struct {
	MapStore domain.MapStore
}

func (handler NewMap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	mapData := r.FormValue("mapData")
	if mapData == "" {
		http.Error(w, "No mapData received", http.StatusBadRequest)
		return
	}
	newMap, err := mapFromData(mapData)
	if err != nil {
		log.Printf("Failed to create map from data: %v\n", err)
		http.Error(w, "Failed to create map from data.", http.StatusInternalServerError)
		return
	}
	err = handler.MapStore.Insert(newMap)
	if err != nil {
		log.Printf("Failed to insert new map into store: %v\n", err)
		http.Error(w, "Failed to insert new map into store.", http.StatusInternalServerError)
		return
	}
	// TODO: redirect (to new challenge page for this map?)
}

func mapFromData(mapData string) (domain.Map, error) {
	newMap := domain.Map{
		MapID:         "",
		Name:          "", // TODO: not yet implemented
		Polygon:       "", // TODO: not yet sent from client
		Area:          -1, // TODO: not yet implemented
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

type Map struct {
	MapStore domain.MapStore
}

func (handler Map) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mapID, ok := r.URL.Query()["id"]
	if !ok || len(mapID) == 0 {
		http.Error(w, "no id!", http.StatusBadRequest)
		log.Printf("no map id!\n")
		return
	}
	foundMap, err := handler.MapStore.Get(mapID[0])
	if err != nil {
		http.Error(w, "failed to get map", http.StatusInternalServerError)
		log.Printf("Failed to get map: %v\n", err)
	}
	json.NewEncoder(w).Encode(foundMap)
	//http.Error(w, "not implemented", http.StatusNotImplemented)
}

type NewChallenge struct {
	ChallengeStore domain.ChallengeStore
}

func (handler *NewChallenge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	challengeData := r.FormValue("challengeData")
	if challengeData == "" {
		http.Error(w, "No challengeData received", http.StatusBadRequest)
		return
	}
	newChallenge, err := challengeFromData(challengeData)
	if err != nil {
		log.Printf("Failed to create challenge from data: %v\n", err)
		http.Error(w, "Failed to create challenge from data.", http.StatusInternalServerError)
		return
	}
	err = handler.ChallengeStore.Insert(newChallenge)
	if err != nil {
		log.Printf("Failed to insert new challenge into store: %v\n", err)
		http.Error(w, "Failed to insert new challenge into store.", http.StatusInternalServerError)
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
		Places:      make([]domain.ChallengePlace, 0),
	}
	var locations []jsonPoint
	if err := json.Unmarshal([]byte(challengeData), &locations); err != nil {
		return newChallenge, err
	}
	// convert from degrees to radians (ffs) and populate challenge.Places
	for i := range locations {
		newChallenge.Places = append(newChallenge.Places, domain.ChallengePlace{
			ChallengeID: newChallenge.ChallengeID,
			Location:    s2.LatLngFromDegrees(locations[i].Lat, locations[i].Lng),
		})
	}

	return newChallenge, nil
}

type NewChallengeResult struct {
	ChallengeResultStore domain.ChallengeResultStore
}

func (handler *NewChallengeResult) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: create new ChallengeResult from form
	// TODO: validate ChallengeResult
	//       nickname not empty
	// TODO: insert into store
	// TODO: redirect to actual game
}
