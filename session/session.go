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

func getSession(r *http.Request) (*sessions.Session, error) {

	store := sessions.NewCookieStore([]byte(os.Getenv("SPOTIFY_SESSION_SECRET")))
	session, err := store.Get(r, "spotify-helper-session")

	session.Options = &sessions.Options{
		MaxAge: 0, // Session
		Path:   "/",
	}

	return session, err
}

func Read(r *http.Request, key string) (value string, err error) {

	session, err := getSession(r)
	if err != nil {
		return "", err
	}

	if session.Values[key] == nil {
		session.Values[key] = ""
	}

	return session.Values[key].(string), err
}

func Write(w http.ResponseWriter, r *http.Request, name string, value string) (err error) {

	session, err := getSession(r)
	if err != nil {
		return err
	}

	session.Values[name] = value

	return session.Save(r, w)
}

func WriteMany(w http.ResponseWriter, r *http.Request, values map[string]string) (err error) {

	session, err := getSession(r)
	if err != nil {
		return err
	}

	for k, v := range values {
		session.Values[k] = v
	}

	return session.Save(r, w)
}

func Clear(w http.ResponseWriter, r *http.Request) (err error) {

	session, err := getSession(r)
	session.Values = make(map[interface{}]interface{})

	if err != nil {
		return err
	}

	err = session.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}

func GetFlashes(w http.ResponseWriter, r *http.Request) (flashes []interface{}, err error) {

	session, err := getSession(r)
	if err != nil {
		return nil, err
	}

	flashes = session.Flashes()
	err = session.Save(r, w)
	if err != nil {
		return nil, err
	}

	return flashes, nil
}

func SetFlash(w http.ResponseWriter, r *http.Request, flash string) (err error) {

	session, err := getSession(r)

	session.AddFlash(flash)

	return session.Save(r, w)
}

func IsLoggedIn(r *http.Request) (val bool, err error) {
	read, err := Read(r, TokenToken)
	return read != "", err
}
