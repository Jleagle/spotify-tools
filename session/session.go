package session

import (
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	Token    = "token"
	LastPage = "last-page"
)

func getCookieStore() *sessions.CookieStore {

	return sessions.NewCookieStore([]byte("something-very-secret"))
}

func Read(w http.ResponseWriter, r *http.Request, key string) string {

	store := getCookieStore()
	session, err := store.Get(r, "month")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if session.Values[key] == nil {
		session.Values[key] = ""
	}

	return session.Values[key].(string)
}

func Write(w http.ResponseWriter, r *http.Request, values map[string]string) (err error) {

	store := getCookieStore()
	session, err := store.Get(r, "month")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	for k, v := range values {
		session.Values[k] = v
	}

	session.Save(r, w)
	return err
}

func IsLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	return Read(w, r, Token) != ""
}

func GetFlashes(w http.ResponseWriter, r *http.Request) []interface{} {

	store := getCookieStore()
	session, err := store.Get(r, "month")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	flashes := session.Flashes()
	session.Save(r, w)
	return flashes
}

func SetFlash(w http.ResponseWriter, r *http.Request, flash string) {

	store := getCookieStore()
	session, err := store.Get(r, "month")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	session.AddFlash(flash)
	session.Save(r, w)
}
