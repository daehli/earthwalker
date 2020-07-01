package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gitlab.com/glatteis/earthwalker/domain"
)

type Challenges struct {
	ChallengeStore domain.ChallengeStore
}

func (handler Challenges) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		challengeID, _ := shiftPath(r.URL.Path)
		if len(challengeID) == 0 || challengeID == "/" {
			sendError(w, "missing challenge id", http.StatusBadRequest)
			return
		}
		foundChallenge, err := handler.ChallengeStore.Get(challengeID)
		if err != nil {
			sendError(w, "failed to get challenge from store", http.StatusInternalServerError)
			log.Printf("Failed to get challenge from store: %v\n", err)
			return
		}
		json.NewEncoder(w).Encode(foundChallenge)
	case http.MethodPost:
		newChallenge, err := challengeFromRequest(r)
		if err != nil {
			sendError(w, "failed to create challenge from request", http.StatusInternalServerError)
			log.Printf("Failed to create challenge from request: %v\n", err)
			return
		}
		err = handler.ChallengeStore.Insert(newChallenge)
		if err != nil {
			sendError(w, "failed to insert challenge into store", http.StatusInternalServerError)
			log.Printf("Failed to insert challenge into store: %v\n", err)
			return
		}
		json.NewEncoder(w).Encode(newChallenge)
	default:
		sendError(w, "api/challenges endpoint does not exist.", http.StatusNotFound)
	}
}

func challengeFromRequest(r *http.Request) (domain.Challenge, error) {
	newChallenge := domain.Challenge{
		Places: make([]domain.ChallengePlace, 0),
	}
	err := json.NewDecoder(r.Body).Decode(&newChallenge)
	if err != nil {
		return newChallenge, fmt.Errorf("failed to decode newChallenge from request: %v", err)
	}
	newChallenge.ChallengeID = domain.RandAlpha(10)
	for i := range newChallenge.Places {
		newChallenge.Places[i].ChallengeID = newChallenge.ChallengeID
	}
	return newChallenge, nil
}
