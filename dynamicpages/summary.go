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
	if err == player.PlayerSessionNotFoundError {
		w.Write([]byte("you are not authenticated to guess!"))
		w.WriteHeader(401)
		return
	} else if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	foundChallenge, err := challenge.GetChallenge(session.GameID)
	if err == challenge.ChallengeNotFoundError {
		w.Write([]byte("this challenge does not exist!"))
		w.WriteHeader(404)
		return
	} else if err != nil {
		log.Println(err)
		w.Write([]byte("there was some kind of internal error, sorry!"))
		w.WriteHeader(500)
		return
	}

	if session.Round() < foundChallenge.Settings.NumRounds {
		w.Write([]byte("You have not completed every round yet, so you cannot view the summary."))
		w.WriteHeader(500)
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
		w.Write([]byte("there was some kind of internal error, sorry!"))
		w.WriteHeader(500)
		return
	}
}
