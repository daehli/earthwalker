package dynamicpages

import (
	"html/template"
	"log"
	"net/http"
)

var getPlaces = template.Must(template.ParseFiles("templates/main_template.html.tmpl", "templates/get_places/get_places.html.tmpl"))

// ServeGetPlaces serves the get_places template.
func ServeGetPlaces(w http.ResponseWriter, r *http.Request) {
	err := getPlaces.Execute(w, struct{}{})
	if err != nil {
		log.Println(err)
		w.Write([]byte("there was some kind of internal error, sorry!"))
		w.WriteHeader(500)
		return
	}
}
