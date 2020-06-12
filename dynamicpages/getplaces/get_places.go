// Package getplaces serves the get_places template.
package getplaces

import (
	"html/template"
	"log"
	"net/http"
	texttemplate "text/template"

	"gitlab.com/glatteis/earthwalker/config"
)

var getPlaces = template.Must(template.ParseFiles(config.Env.EarthwalkerStaticPath+"/templates/main_template.html.tmpl", config.Env.EarthwalkerStaticPath+"/templates/get_places/get_places.html.tmpl"))
var getPlacesJS = texttemplate.Must(texttemplate.ParseFiles(config.Env.EarthwalkerStaticPath + "/templates/get_places/get_places.js.tmpl"))

// ServeGetPlaces serves the get_places template.
func ServeGetPlaces(w http.ResponseWriter, r *http.Request) {
	err := getPlaces.Execute(w, struct{}{})
	if err != nil {
		log.Println(err)
		http.Error(w, "there was some kind of internal error, sorry!", http.StatusUnprocessableEntity)
		return
	}
}

// ServeGetPlacesJS serves the get_places.js template.
func ServeGetPlacesJS(w http.ResponseWriter, r *http.Request) {
	err := getPlacesJS.Execute(w, struct {
		Config config.FileType
	}{
		Config: config.File,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, "there was some kind of internal error, sorry!", http.StatusUnprocessableEntity)
		return
	}
}
