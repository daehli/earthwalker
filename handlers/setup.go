// handlers in this file create and store new structs before the game begins
// (Map, Challenge, ChallengeResult, etc.)
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gitlab.com/glatteis/earthwalker/domain"
)

type NewMap struct {
	MapStore domain.MapStore
}

func (handler NewMap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	newMap, err := mapFromRequest(r)
	if err != nil {
		log.Printf("Failed to create map from request: %v\n", err)
		http.Error(w, "Failed to create map from request.", http.StatusInternalServerError)
		return
	}
	err = handler.MapStore.Insert(newMap)
	if err != nil {
		log.Printf("Failed to insert new map into store: %v\n", err)
		http.Error(w, "Failed to insert new map into store.", http.StatusInternalServerError)
		return
	}
	// TODO: redirect (to new challenge page for this map?)
	// TODO: remove this debugging response
	http.Redirect(w, r, "/map?id="+newMap.MapID, http.StatusFound)
}

func mapFromRequest(r *http.Request) (domain.Map, error) {
	newMap := domain.Map{}
	err := json.NewDecoder(r.Body).Decode(&newMap)
	if err != nil {
		return newMap, fmt.Errorf("Failed to decode newMap from request: %v", err)
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
		return
	}
	json.NewEncoder(w).Encode(foundMap)
}

type NewChallenge struct {
	ChallengeStore domain.ChallengeStore
}

func (handler NewChallenge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	newChallenge, err := challengeFromRequest(r)
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
	// TODO: redirect (to new challenge page for this challenge?)
	// TODO: remove this debugging response
	w.Write([]byte(newChallenge.ChallengeID))
}

func challengeFromRequest(r *http.Request) (domain.Challenge, error) {
	newChallenge := domain.Challenge{
		Places: make([]domain.ChallengePlace, 0),
	}
	err := json.NewDecoder(r.Body).Decode(&newChallenge)
	if err != nil {
		return newChallenge, fmt.Errorf("Failed to decode newChallenge from request: %v", err)
	}
	newChallenge.ChallengeID = domain.RandAlpha(10)
	for _, place := range newChallenge.Places {
		place.ChallengeID = newChallenge.ChallengeID
	}
	return newChallenge, err
}

type Challenge struct {
	ChallengeStore domain.ChallengeStore
}

func (handler Challenge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	challengeID, ok := r.URL.Query()["id"]
	if !ok || len(challengeID) == 0 {
		http.Error(w, "no id!", http.StatusBadRequest)
		log.Printf("no challenge id!\n")
		return
	}
	foundChallenge, err := handler.ChallengeStore.Get(challengeID[0])
	if err != nil {
		http.Error(w, "failed to get challenge", http.StatusInternalServerError)
		log.Printf("Failed to get challenge: %v\n", err)
		return
	}
	json.NewEncoder(w).Encode(foundChallenge)
}

type NewChallengeResult struct {
	ChallengeResultStore domain.ChallengeResultStore
}

func (handler NewChallengeResult) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: create new ChallengeResult from form
	// TODO: validate ChallengeResult
	//       nickname not empty
	// TODO: insert into store
	// TODO: redirect to actual game
}
