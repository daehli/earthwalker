package player

import (
	"errors"
	"log"
	"net/http"
	"time"
)

var PlayerSessionNotFoundError = errors.New("no player session found")

// SetSessionCookie sets the session cookie of a session into the browser.
func SetSessionCookie(session PlayerSession, w http.ResponseWriter) {
	c := http.Cookie{
		Name:   "earthwalker-session",
		Value:  session.UniqueIdentifier,
		MaxAge: int((24 * time.Hour).Seconds()),
	}
	http.SetCookie(w, &c)
}

// GetSessionFromCookie retrieves the cookie from a session
func GetSessionFromCookie(r *http.Request) (PlayerSession, error) {
	var cookie *http.Cookie
	for _, c := range r.Cookies() {
		log.Println(*c)
		if c.Name == "earthwalker-session" {
			log.Println("Found cookie")
			cookie = c
		}
	}
	if cookie == nil {
		return PlayerSession{}, PlayerSessionNotFoundError
	}

	session, err := loadPlayerSession(cookie.Value)
	if err != nil {
		return PlayerSession{}, err
	}

	return session, nil
}
