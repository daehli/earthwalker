package placefinder

import (
	"encoding/json"
	"github.com/golang/geo/s2"
	"gitlab.com/glatteis/earthwalker/challenge"
	"log"
	"net/http"
)

func RespondToPoints(w http.ResponseWriter, r *http.Request) {
	type jsonPoint struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}

	r.ParseForm()
	result := r.FormValue("result")

	var content []jsonPoint

	if err := json.Unmarshal([]byte(result), &content); err != nil {
		http.Error(w, err.Error(), 403)
		return
	}

	locations := make([]s2.LatLng, len(content))

	for i := range content {
		locations[i] = s2.LatLngFromDegrees(content[i].Lat, content[i].Lng)
	}

	resultingChallenge, err := challenge.NewChallenge(locations)

	if err != nil {
		log.Println(err)
		w.Write([]byte("Internal server error! Please contact an administrator or something"))
		return
	}

	http.Redirect(w, r, "/game?c="+resultingChallenge.UniqueIdentifier, 302)
}
