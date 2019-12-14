package dynamicpages

import (
	"github.com/pkg/errors"
	"gitlab.com/glatteis/earthwalker/challenge"
	"gitlab.com/glatteis/earthwalker/player"
	"html/template"
	"log"
	"net/http"
	"sort"
)

var summaryPage = template.Must(template.ParseFiles("templates/main_template.html.tmpl", "templates/summary/summary.html.tmpl"))

type rankingType struct {
	Nickname  string
	NumPoints int
}

type summaryServeStruct struct {
	Rankings []rankingType
}

// ServeSummary serves the summary page.
func ServeSummary(w http.ResponseWriter, r *http.Request) {
	session, err := player.GetSessionFromCookie(r)
	if err == player.ErrPlayerSessionNotFound {
		http.Error(w, "you are not authenticated to guess!", http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "some error occured", http.StatusUnprocessableEntity)
		return
	}

	foundChallenge, err := challenge.GetChallenge(session.GameID)
	if err == challenge.ErrChallengeNotFound {
		http.Error(w, "this challenge does not exist!", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "there was some kind of internal error, sorry!", http.StatusUnprocessableEntity)
		return
	}

	if session.Round() < foundChallenge.Settings.NumRounds {
		http.Error(w, "You have not completed every round yet, so you cannot view the summary.", http.StatusUnprocessableEntity)
		return
	}

	ranking := make([]rankingType, 0)

	for _, playerThatCompleted := range foundChallenge.Guesses[len(foundChallenge.Guesses)-1] {
		completedSession, err := player.LoadPlayerSession(playerThatCompleted.PlayerID)
		if err != nil {
			log.Println(errors.Wrap(err, "while loading a player that should have guessed"))
		}
		var sumPoints int
		for _, p := range completedSession.Points {
			sumPoints += p
		}
		ranking = append(ranking, rankingType{
			Nickname:  completedSession.Nickname,
			NumPoints: sumPoints,
		})
	}

	sort.Slice(ranking, func(i, j int) bool {
		return ranking[i].NumPoints >= ranking[j].NumPoints
	})

	err = summaryPage.Execute(w, summaryServeStruct{
		Rankings: ranking,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, "there was some kind of internal error, sorry!", http.StatusUnprocessableEntity)
		return
	}
}
