package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gitlab.com/glatteis/earthwalker/domain"
)

type Results struct {
	ChallengeResultStore domain.ChallengeResultStore
}

func (handler Results) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		challengeResultID, _ := shiftPath(r.URL.Path)
		if len(challengeResultID) == 0 || challengeResultID == "/" {
			sendError(w, "missing result id", http.StatusBadRequest)
			return
		}
		foundChallengeResult, err := handler.ChallengeResultStore.Get(challengeResultID)
		if err != nil {
			sendError(w, "failed to get result from store", http.StatusInternalServerError)
			log.Printf("Failed to get result from store: %v\n", err)
			return
		}
		json.NewEncoder(w).Encode(foundChallengeResult)
	case http.MethodPost:
		newChallengeResult, err := challengeResultFromRequest(r)
		if err != nil {
			sendError(w, "failed to create result from request", http.StatusInternalServerError)
			log.Printf("Failed to create result from request: %v\n", err)
			return
		}
		err = handler.ChallengeResultStore.Insert(newChallengeResult)
		if err != nil {
			sendError(w, "failed to insert result into store", http.StatusInternalServerError)
			log.Printf("Failed to insert result into store: %v\n", err)
			return
		}
		json.NewEncoder(w).Encode(newChallengeResult)
	default:
		sendError(w, "api/results endpoint does not exist.", http.StatusNotFound)
	}
}

func challengeResultFromRequest(r *http.Request) (domain.ChallengeResult, error) {
	newChallengeResult := domain.ChallengeResult{
		Guesses: make([]domain.Guess, 0),
	}
	err := json.NewDecoder(r.Body).Decode(&newChallengeResult)
	if err != nil {
		return newChallengeResult, fmt.Errorf("failed to decode newChallengeResult from request: %v", err)
	}
	newChallengeResult.ChallengeResultID = domain.RandAlpha(10)
	return newChallengeResult, nil
}
