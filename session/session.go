package session

import (
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	LastPage    = "last-page"
	AuthToken   = "auth.token"
	AuthState   = "auth.state"
	UserCountry = "user.country"
	UserID      = "user.id"
)

func getSession(r *http.Request) *sessions.Session {

	store := sessions.NewCookieStore([]byte("something-very-secret"))
	session, err := store.Get(r, "hour")
	if err != nil {
		// todo, show error page
		//http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	session.Options = &sessions.Options{
		MaxAge: 0, // Session
		Path:   "/",
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
	return Read(r, AuthToken) != ""
}
