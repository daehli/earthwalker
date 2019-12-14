// Package urlbuilder builds google street view urls from coordinates. this may
// break at any time, because it uses undocumented parameters.
package urlbuilder

import (
	"github.com/golang/geo/s2"
	"log"
	"net/url"
	"strconv"
)

func floatToString(number float64) string {
	return strconv.FormatFloat(number, 'f', 14, 64)
}

// BuildURL builds google street view urls from coordinates
func BuildURL(location s2.LatLng) string {
	baseURL, err := url.Parse("https://www.google.com/maps")
	if err != nil {
		log.Fatal("Failed while parsing static gmaps url", err)
	}
	query := baseURL.Query()
	// see https://stackoverflow.com/questions/387942/google-street-view-url
	// for a reverse-engineering of the parameters

	// the layer must be set to c (the street view layer)
	query.Set("layer", "c")
	// latitude and longitude go into parameter cbll
	query.Set("cbll", floatToString(location.Lat.Degrees())+","+floatToString(location.Lng.Degrees()))

	baseURL.RawQuery = query.Encode()

	return baseURL.String()
}
