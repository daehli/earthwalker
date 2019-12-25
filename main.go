// earthwalker © 2019 Linus Heck

// earthwalker is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Package main is the main package of earthwalker.
package main

import (
	"flag"
	"gitlab.com/glatteis/earthwalker/challenge"
	"gitlab.com/glatteis/earthwalker/database"
	"gitlab.com/glatteis/earthwalker/dynamicpages/continuegame"
	"gitlab.com/glatteis/earthwalker/dynamicpages/getplaces"
	"gitlab.com/glatteis/earthwalker/dynamicpages/modifyfrontend"
	"gitlab.com/glatteis/earthwalker/dynamicpages/scorepage"
	"gitlab.com/glatteis/earthwalker/dynamicpages/setnickname"
	"gitlab.com/glatteis/earthwalker/dynamicpages/summary"
	"gitlab.com/glatteis/earthwalker/placefinder"
	"gitlab.com/glatteis/earthwalker/player"
	"gitlab.com/glatteis/earthwalker/streetviewserver"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func main() {
	defer database.CloseDB()

	rand.Seed(time.Now().UnixNano())
	port := flag.Int("port", 8080, "the port the server is running on")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		session, err := player.GetSessionFromCookie(r)
		if err != nil || session.GameID == "" {
			getplaces.ServeGetPlaces(w, r)
			return
		}
		continuegame.ServeContinueGame(w, r)
	})
	http.HandleFunc("/continue", func(w http.ResponseWriter, r *http.Request) {
		session, err := player.GetSessionFromCookie(r)
		if err != nil {
			getplaces.ServeGetPlaces(w, r)
			return
		}
		redirectURL := "/game?c=" + session.GameID
		http.Redirect(w, r, redirectURL, http.StatusFound)
	})
	http.HandleFunc("/newgame", getplaces.ServeGetPlaces)
	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		challenge.ServeChallenge(w, r)
	})
	http.HandleFunc("/maps/", streetviewserver.ServeMaps)

	http.HandleFunc("/found_points", placefinder.RespondToPoints)

	http.HandleFunc("/scores", scorepage.ServeScores)
	http.HandleFunc("/set_nickname", setnickname.ServeSetNickname)
	http.HandleFunc("/summary", summary.ServeSummary)
	http.HandleFunc("/modify_frontend.js", modifyfrontend.ServeModifyFrontend)

	http.HandleFunc("/guess", challenge.HandleGuess)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("earthwalker is running on ", *port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
