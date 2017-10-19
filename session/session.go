package session

import (
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	Token    = "token"     // Spotify access token
	LastPage = "last-page" // The last page a user was on for page flow
	State    = "state"     // For OAuth security
)

func getSession(r *http.Request) *sessions.Session {

	store := sessions.NewCookieStore([]byte("something-very-secret"))
	session, err := store.Get(r, "month")
	if err != nil {
		// todo, show error page
		//http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return session
}

func Read(r *http.Request, key string) string {

	session := getSession(r)

	if session.Values[key] == nil {
		session.Values[key] = ""
	}

	return session.Values[key].(string)
}

func Write(w http.ResponseWriter, r *http.Request, name string, value string) {

	session := getSession(r)
	session.Values[name] = value

	err := session.Save(r, w)
	if err != nil {
		// todo, show error page
	}
}

func GetFlashes(w http.ResponseWriter, r *http.Request) []interface{} {

	session := getSession(r)

	flashes := session.Flashes()
	err := session.Save(r, w)
	if err != nil {
		// todo, show error page
	}
	return flashes
}

func SetFlash(w http.ResponseWriter, r *http.Request, flash string) {

	session := getSession(r)

	session.AddFlash(flash)
	err := session.Save(r, w)
	if err != nil {
		// todo, show error page
	}
}

func IsLoggedIn(r *http.Request) bool {
	return Read(r, Token) != ""
}
