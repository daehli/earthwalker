package handlers

import (
	"regexp"
)

// Eliminate street names from packets
// match [["Jl. SMA Aek Kota Batu","id"],["Sumatera Utara","de"]] (see regex_test.go)
var streetNameChars = "(\\p{L}| |\\d|\\_|\\-|\\,|\\.|/)"
var listOfStreetNames = "\\[\"" + streetNameChars + "+\"+,\"" + streetNameChars + "{1,10}\"\\]"
var streetNameRegex *regexp.Regexp = regexp.MustCompile(listOfStreetNames)

// Eliminate shop icons from packets. These are an object in json. for an example, see regex_test.go.
// Heuristic used: Always starts with [[[ and ends with ]]], not followed by a [\".
// Always contains https://maps.gstatic.com/mapfiles/annotations/icons/ somewhere.
// var shopChars = "(\\p{L}| |\\d|\\_|\\-|\\,|\\.|\\[|\\])"
// var shopIcons = "[[[" + shopChars + "*?" + "https:\\/\\/maps.gstatic.com\\/mapfiles\\/annotations\\/icons\\/" + shopChars + "*?" + "]]],[^\\[\"]"
var shopIcons = "\\[\\[\\[([\\[\\],\":\\/. \\\\_\\-\\pL\\pN\\pM\\pZ\\pP\\pC])*https:\\/\\/maps\\.gstatic\\.com\\/mapfiles\\/annotations\\/icons\\/([0-9\\[\\],\":\\/. \\\\_\\-\\pL])*\\]\\]\\],"
var shopIconRegex *regexp.Regexp = regexp.MustCompile(shopIcons)

// filterStrings filters all string contents from a given string (as byte array),
// used to strip all localization information from a specific street view packet
func filterStrings(body []byte) []byte {
	result := streetNameRegex.ReplaceAllString(string(body), "[\"\",\"\"]")
	result = shopIconRegex.ReplaceAllString(result, "[],")
	return []byte(result)
}
