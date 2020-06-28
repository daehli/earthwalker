package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gitlab.com/glatteis/earthwalker/domain"
)

type Guesses struct {
	ChallengeResultStore domain.ChallengeResultStore
}

func (handler Guesses) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		log.Println("handling guess POST!")
		newGuess, err := guessFromRequest(r)
		if err != nil {
			sendError(w, "failed to create guess from request", http.StatusInternalServerError)
			log.Printf("Failed to create guess from request: %v\n", err)
			return
		}
		result, err := handler.ChallengeResultStore.Get(newGuess.ChallengeResultID)
		if err != nil {
			sendError(w, "failed to get result specified in guess", http.StatusInternalServerError)
			log.Printf("Failed to get result specified in guess: %v\n", err)
			return
		}
		if len(result.Guesses) != newGuess.RoundNum {
			sendError(w, "guess round num does not match existing result", http.StatusInternalServerError)
			log.Printf("Guess round num does not match existing result: %v\n", err)
			return
		}
		result.Guesses = append(result.Guesses, newGuess)
		err = handler.ChallengeResultStore.Insert(result)
		if err != nil {
			sendError(w, "failed to insert guess into store", http.StatusInternalServerError)
			log.Printf("Failed to insert result with new guess into store: %v\n", err)
			return
		}
		json.NewEncoder(w).Encode(result)
	default:
		sendError(w, "api/guesses endpoint does not exist.", http.StatusNotFound)
	}
}

// TODO: FIXME: currently taking latlng in degrees from the client and storing that number as "radians"
//              doesn't matter at the moment because the server doesn't do anything with guesses
func guessFromRequest(r *http.Request) (domain.Guess, error) {
	newGuess := domain.Guess{}
	err := json.NewDecoder(r.Body).Decode(&newGuess)
	if err != nil {
		return newGuess, fmt.Errorf("failed to decode newGuess from request: %v", err)
	}
	return newGuess, nil
}
