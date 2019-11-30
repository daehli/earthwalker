package main

import (
	"gitlab.com/glatteis/earthwalker/placefinder"
	"gitlab.com/glatteis/earthwalker/streetviewserver"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/get_places/get_places.html", 301)
	})
	http.HandleFunc("/maps/", streetviewserver.ServeMaps)

	http.HandleFunc("/found_points", placefinder.RespondToPoints)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
