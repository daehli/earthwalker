package handlers

// handlers in this file serve the frontend

import (
	htemplate "html/template"
	"log"
	"net/http"
	ttemplate "text/template"
)

type Root struct{}

func (handler Root) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	challengeID, err := challengeIDFromRequest(r)
	if err != nil || challengeID == "" {
		http.Redirect(w, r, "/createmap", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/play", http.StatusFound)
}

type DynamicHTML struct {
	Template *htemplate.Template
	Data     interface{}
}

func (handler DynamicHTML) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := handler.Template.Execute(w, handler.Data)
	if err != nil {
		log.Printf("Failed to serve html template: %v\n", err)
		http.Error(w, "Failed to serve page.", http.StatusInternalServerError)
	}
}

type DynamicText struct {
	Template *ttemplate.Template
	Data     interface{}
}

func (handler DynamicText) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := handler.Template.Execute(w, handler.Data)
	if err != nil {
		log.Printf("Failed to serve text template: %v\n", err)
		http.Error(w, "Failed to serve file.", http.StatusInternalServerError)
	}
}
