package main

import (
	"fmt"
	"net/http"

	"github.com/Jleagle/go-helpers/helpers"
	"github.com/Jleagle/spotifyhelper/session"
	spot "github.com/Jleagle/spotifyhelper/spotify"
	"github.com/Jleagle/spotifyhelper/structs"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {

	if session.IsLoggedIn(r) {
		postlogin(w, r)
		return
	}

	if r.URL.Query().Get("auth") == "1" {
		state := helpers.RandomString(6, "abcdefghijklmnopqrstuvwxyz")
		session.Write(w, r, session.AuthState, state)

		auth := spot.GetAuthenticator()
		http.Redirect(w, r, auth.AuthURL(state), 302)
		return
	}

	returnTemplate(w, r, "login", structs.TemplateVars{}, nil)
	return
}

func loginCallbackHandler(w http.ResponseWriter, r *http.Request) {

	if session.IsLoggedIn(r) {
		postlogin(w, r)
		return
	}

	queryErr := r.URL.Query().Get("error")
	if queryErr != "" {

		vars := structs.TemplateVars{}
		vars.ErrorMessage = "Spotify: " + queryErr

		returnTemplate(w, r, "error", vars, nil)
		return
	}

	auth := spot.GetAuthenticator()
	state := session.Read(r, session.AuthState)

	tok, err := auth.Token(state, r)
	if err != nil {
		returnTemplate(w, r, "error", structs.TemplateVars{}, err)
		return
	}

	client := auth.NewClient(tok)
	user, err := client.CurrentUser()
	if err != nil {
		fmt.Println("Getting user details: " + err.Error())
	}

	// todo, grab stuff from user, save to db?
	session.Write(w, r, session.AuthState, "")
	session.Write(w, r, session.AuthToken, tok.AccessToken)
	session.Write(w, r, session.UserCountry, user.Country)
	session.Write(w, r, session.UserID, user.ID)

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

	// todo, do this in one call
	session.Write(w, r, session.AuthToken, "")
	session.Write(w, r, session.AuthState, "")
	session.Write(w, r, session.LastPage, "")
	session.Write(w, r, session.UserCountry, "")
	session.Write(w, r, session.UserID, "")

	session.SetFlash(w, r, "Logged Out!")

	http.Redirect(w, r, lastPage, 302)
}
