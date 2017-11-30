package session

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

const (
	LastPage     = "last-page"
	AuthState    = "auth.state"
	UserCountry  = "user.country"
	UserID       = "user.id"
	TokenToken   = "token.token"
	TokenType    = "token.type"
	TokenRefresh = "token.refresh"
	TokenExpiry  = "token.expiry"
)

func getSession(r *http.Request) *sessions.Session {

	store := sessions.NewCookieStore([]byte(os.Getenv("SPOTIFY_SESSION_SECRET")))
	session, err := store.Get(r, "spotify-helper-session")
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

func WriteMany(w http.ResponseWriter, r *http.Request, values map[string]string) {

	session := getSession(r)
	for k, v := range values {
		session.Values[k] = v
	}

	err := session.Save(r, w)
	if err != nil {
		// todo, show error page
	}
}

func Clear(w http.ResponseWriter, r *http.Request) {

	session := getSession(r)
	session.Values = make(map[interface{}]interface{})

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
	return Read(r, TokenToken) != ""
}
