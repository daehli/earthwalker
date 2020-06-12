package player

import (
	"errors"
	"net/http"
	"time"
)

// ErrPlayerSessionNotFound is the error that occurs when no player session is found.
var ErrPlayerSessionNotFound = errors.New("no player session found")

// SetSessionCookie sets the session cookie of a session into the browser.
func SetSessionCookie(session Session, w http.ResponseWriter) {
	c := http.Cookie{
		Name:   "earthwalker-session",
		Value:  session.UniqueIdentifier,
		MaxAge: int((24 * time.Hour).Seconds()),
	}
	http.SetCookie(w, &c)
}

// GetSessionFromCookie retrieves the cookie from a session
func GetSessionFromCookie(r *http.Request) (Session, error) {
	var cookie *http.Cookie
	for _, c := range r.Cookies() {
		if c.Name == "earthwalker-session" {
			cookie = c
		}
	}
	if cookie == nil {
		return Session{}, ErrPlayerSessionNotFound
	}

	session, err := LoadPlayerSession(cookie.Value)
	if err != nil {
		return Session{}, err
	}

	return session, nil
}
