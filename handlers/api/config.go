package api

import (
	"fmt"
	"log"
	"net/http"

	"gitlab.com/glatteis/earthwalker/domain"
)

type Config struct {
	Config domain.Config
}

func (handler Config) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("api/config accepts only GET requests, not '%s'.", r.Method), http.StatusNotFound)
	}
	var respJSON string
	switch head, _ := shiftPath(r.URL.Path); head {
	case "tileserver":
		respJSON = "{tileserver: \"" + handler.Config.TileServerURL + "\"}"
	case "nolabeltileserver":
		respJSON = "{nolabeltileserver: \"" + handler.Config.NoLabelTileServerURL + "\"}"
	default:
		sendError(w, fmt.Sprintf("api/config endpoint '%s' does not exist.", r.URL.Path), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(respJSON))
	if err != nil {
		log.Printf("Error writing response: %v\n", err)
	}
}
