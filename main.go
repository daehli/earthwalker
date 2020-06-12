// earthwalker © 2019-2020 Linus Heck & Contributors

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
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"gitlab.com/glatteis/earthwalker/challenge"
	"gitlab.com/glatteis/earthwalker/config"
	"gitlab.com/glatteis/earthwalker/database"
	"gitlab.com/glatteis/earthwalker/dynamicpages/beforestart"
	"gitlab.com/glatteis/earthwalker/dynamicpages/continuegame"
	"gitlab.com/glatteis/earthwalker/dynamicpages/getplaces"
	"gitlab.com/glatteis/earthwalker/dynamicpages/modifyfrontend"
	"gitlab.com/glatteis/earthwalker/dynamicpages/scorepage"
	"gitlab.com/glatteis/earthwalker/dynamicpages/setnickname"
	"gitlab.com/glatteis/earthwalker/dynamicpages/summary"
	"gitlab.com/glatteis/earthwalker/placefinder"
	"gitlab.com/glatteis/earthwalker/player"
	"gitlab.com/glatteis/earthwalker/streetviewserver"
)

var placesAndFunctions = map[string]func(w http.ResponseWriter, r *http.Request){
	"/newgame":            getplaces.ServeGetPlaces,
	"/beforestart":        beforestart.ServeBeforeStart,
	"/game":               challenge.ServeChallenge,
	"/maps/":              streetviewserver.ServeMaps,
	"/found_points":       placefinder.RespondToPoints,
	"/scores":             scorepage.ServeScores,
	"/set_nickname":       setnickname.ServeSetNickname,
	"/summary":            summary.ServeSummary,
	"/modify_frontend.js": modifyfrontend.ServeModifyFrontend,
	"/guess":              challenge.HandleGuess,
}

func cleanup() {
	database.CloseDB()
}

func main() {
	// Either defer cleanup for when the program exits...
	defer cleanup()
	// Or listen for SIGTERM and also clean up.
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(0)
	}()

	rand.Seed(time.Now().UnixNano())

	// get port from config
	port := config.Env.EarthwalkerPort
	if port == "" {
		portFlag := flag.Int("port", 8080, "the port the server is running on")
		flag.Parse()
		port = strconv.Itoa(*portFlag)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		session, err := player.GetSessionFromCookie(r)
		if err != nil || session.GameID == "" {
			getplaces.ServeGetPlaces(w, r)
			return
		}
		continuegame.ServeContinueGame(w, r, session.Nickname)
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
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(config.Env.EarthwalkerStaticPath+"/static"))))
	for path, function := range placesAndFunctions {
		http.HandleFunc(path, function)
	}

	log.Println("earthwalker is running on ", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
