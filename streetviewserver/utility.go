package streetviewserver

import (
	"regexp"
)

// GOAL : MATCH [["Jl. SMA Aek Kota Batu","id"],["Sumatera Utara","de"]]
var stringRegex = "(\\p{L}| |\\d|\\_|\\-|\\,|\\.|/)"
var languageRegex = "\\[\"" + stringRegex + "+\"+,\"" + stringRegex + "{1,10}\"\\]"

var compiledRegexp *regexp.Regexp = regexp.MustCompile(languageRegex)

// filterStrings filters all string contents from a given string (as byte array),
// used to strip all localization information from a specific street view packet
func filterStrings(body []byte) []byte {
	result := compiledRegexp.ReplaceAllString(string(body), "[\"\",\"\"]")
	return []byte(result)
}
