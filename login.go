package main

import (
	"net/http"
	"os"

	"github.com/Jleagle/spotify-tools/helpers"
	"github.com/Jleagle/spotify-tools/session"
	"github.com/zmb3/spotify"
)

func getAuthenticator() (authenticator spotify.Authenticator) {
	authenticator = spotify.NewAuthenticator("http://localhost:8084/login-callback", spotify.ScopeUserReadPrivate)
	authenticator.SetAuthInfo(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	return authenticator
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	if session.IsLoggedIn(r) {
		postlogin(w, r)
		return
	}

	if r.URL.Query().Get("auth") == "1" {
		state := helpers.RandomString(6)
		session.Write(w, r, session.State, state)

		auth := getAuthenticator()
		http.Redirect(w, r, auth.AuthURL(state), 302)
		return
	}

	returnTemplate(w, r, "login", templateVars{}, nil)
	return
}

func loginCallbackHandler(w http.ResponseWriter, r *http.Request) {

	if session.IsLoggedIn(r) {
		postlogin(w, r)
		return
	}

	queryErr := r.URL.Query().Get("error")
	if queryErr != "" {

		vars := templateVars{}
		vars.ErrorMessage = "Spotify: " + queryErr

		returnTemplate(w, r, "error", vars, nil)
		return
	}

	auth := getAuthenticator()
	state := session.Read(r, session.State)

	tok, err := auth.Token(state, r)
	if err != nil {
		returnTemplate(w, r, "error", templateVars{}, err)
		return
	}

	session.Write(w, r, session.State, "")
	session.Write(w, r, session.Token, tok.AccessToken)

	postlogin(w, r)
	return
}

func postlogin(w http.ResponseWriter, r *http.Request) {

	session.SetFlash(w, r, "Logged In :)")
	lastPage := session.Read(r, session.LastPage)
	http.Redirect(w, r, lastPage, 302)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {

	session.Write(w, r, session.Token, "")
	session.SetFlash(w, r, "Logged Out!")

	http.Redirect(w, r, "/", 302)
}
