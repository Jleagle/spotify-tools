package main

import (
	"net/http"

	"github.com/Jleagle/go-helpers/helpers"
	"github.com/Jleagle/spotify-tools/session"
	spot "github.com/Jleagle/spotify-tools/spotify"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {

	if session.IsLoggedIn(r) {
		postlogin(w, r)
		return
	}

	if r.URL.Query().Get("auth") == "1" {
		state := helpers.RandomString(6, "abcdefghijklmnopqrstuvwxyz")
		session.Write(w, r, session.State, state)

		auth := spot.GetAuthenticator()
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

	auth := spot.GetAuthenticator()
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
	return
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {

	lastPage := session.Read(r, session.LastPage)

	session.Write(w, r, session.Token, "")
	session.Write(w, r, session.State, "")
	session.Write(w, r, session.LastPage, "")

	session.SetFlash(w, r, "Logged Out!")

	http.Redirect(w, r, lastPage, 302)
}
