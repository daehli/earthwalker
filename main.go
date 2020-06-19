// earthwalker © 2019-2020 Linus Heck & Contributors

// earthwalker is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Package main is the main package of earthwalker.
package main

import (
	"flag"
	htemplate "html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	ttemplate "text/template"
	"time"

	"gitlab.com/glatteis/earthwalker/badgerdb"
	"gitlab.com/glatteis/earthwalker/config"
	"gitlab.com/glatteis/earthwalker/handlers"
)

func main() {
	// TODO: can we get rid of this?
	rand.Seed(time.Now().UnixNano())

	// == CONFIG ========
	conf, err := config.Read()
	if err != nil {
		log.Fatalf("Failed to read config: %v\n", err)
	}

	// get port from flag
	// TODO: can we get rid of this?
	port := conf.Port
	if port == "" {
		portFlag := flag.Int("port", 8080, "the port the server is running on")
		flag.Parse()
		port = strconv.Itoa(*portFlag)
	}

	// == DATABASE ========
	db, err := badgerdb.Init(conf.DBPath)
	if err != nil {
		log.Fatalf("Failed to open db at %s: %v\n", conf.DBPath, err)
	}

	// Either defer cleanup for when the program exits...
	defer badgerdb.Close(db)
	// Or listen for SIGTERM and also clean up.
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		badgerdb.Close(db)
		os.Exit(0)
	}()

	// == HANDLERS ========
	var mainTemplate string = conf.StaticPath + "/templates/main_template.html.tmpl"
	http.Handle("/", handlers.Root{})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(conf.StaticPath+"/static"))))
	http.Handle("/mapeditor", handlers.DynamicHTML{Template: htemplate.Must(htemplate.ParseFiles(mainTemplate, conf.StaticPath+"/templates/mapeditor.html.tmpl")), Data: conf})
	http.Handle("/mapeditor.js", handlers.DynamicText{Template: ttemplate.Must(ttemplate.ParseFiles(conf.StaticPath + "/templates/mapeditor.js.tmpl")), Data: conf})
	http.Handle("/get_places.js", handlers.DynamicText{Template: ttemplate.Must(ttemplate.ParseFiles(conf.StaticPath + "/templates/get_places.js.tmpl")), Data: conf})
	http.HandleFunc("/map?id=", func(w http.ResponseWriter, r *http.Request) {
		// serve existing map json
	})

	// == ENGAGE ========
	log.Println("earthwalker is running on ", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
