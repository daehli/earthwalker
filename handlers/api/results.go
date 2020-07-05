package api

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"gitlab.com/glatteis/earthwalker/domain"
)

type Results struct {
	ChallengeResultStore domain.ChallengeResultStore
}

func (handler Results) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	head, tail := shiftPath(r.URL.Path)
	switch head {
	case "all":
		switch r.Method {
		case http.MethodGet:
			challengeID, _ := shiftPath(tail)
			if len(challengeID) == 0 || challengeID == "/" {
				sendError(w, "missing challenge id", http.StatusBadRequest)
				return
			}
			foundChallengeResults, err := handler.ChallengeResultStore.GetAll(challengeID)
			if err != nil {
				sendError(w, "failed to get results from store", http.StatusInternalServerError)
				log.Printf("Failed to get results from store: %v\n", err)
				return
			}
			json.NewEncoder(w).Encode(foundChallengeResults)
		default:
			sendError(w, "api/results/all endpoint does not exist.", http.StatusNotFound)
		}
	default:
		switch r.Method {
		case http.MethodGet:
			challengeResultID := head
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
			// TODO: results don't seem to be echoing as expected?
			json.NewEncoder(w).Encode(newChallengeResult)
		default:
			sendError(w, "api/results endpoint does not exist.", http.StatusNotFound)
		}
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
	newChallengeResult.Icon = rand.Intn(199) + 1
	return newChallengeResult, nil
}
