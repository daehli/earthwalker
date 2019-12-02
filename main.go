package main

import (
	"flag"
	"gitlab.com/glatteis/earthwalker/challenge"
	"gitlab.com/glatteis/earthwalker/database"
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
		log.Println(session, err)
		if err != nil || session.CurrentGameID == "" {
			log.Println(err)
			http.Redirect(w, r, "/static/get_places/get_places.html", 302)
		} else {
			redirectURL := "/game?c=" + session.CurrentGameID
			if session.CurrentRound != 0 {
				redirectURL += "&round=" + strconv.Itoa(session.CurrentRound)
			}
			http.Redirect(w, r, redirectURL, 302)
		}
	})
	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		challenge.ServeChallenge(w, r)
	})
	http.HandleFunc("/maps/", streetviewserver.ServeMaps)

	http.HandleFunc("/found_points", placefinder.RespondToPoints)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
