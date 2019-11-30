package main

import (
	"flag"
	"gitlab.com/glatteis/earthwalker/placefinder"
	"gitlab.com/glatteis/earthwalker/streetviewserver"
	"log"
	"net/http"
	"strconv"
)

func main() {
	port := flag.Int("port", 8080, "the port the server is running on")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/get_places/get_places.html", 301)
	})
	http.HandleFunc("/maps/", streetviewserver.ServeMaps)

	http.HandleFunc("/found_points", placefinder.RespondToPoints)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
