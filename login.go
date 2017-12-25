package main

import (
	"net/http"
	"strconv"

	"github.com/Jleagle/go-helpers/helpers"
	"github.com/Jleagle/go-helpers/rollbar"
	"github.com/Jleagle/spotifyhelper/session"
	"github.com/Jleagle/spotifyhelper/spotify"
	"github.com/Jleagle/spotifyhelper/structs"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {

	vars := structs.TemplateVars{}

	// Check if logged in
	loggedIn, err := session.IsLoggedIn(r)
	if err != nil {
		returnTemplate(w, r, "error", vars, err)
		return
	}
	if loggedIn {
		postlogin(w, r)
		return
	}

	// Create random state
	state := helpers.RandomString(6, "abcdefghijklmnopqrstuvwxyz")
	err = session.Write(w, r, session.AuthState, state)
	if err != nil {
		returnTemplate(w, r, "error", vars, err)
		return
	}

	auth := spotify.GetAuthenticator(r)
	http.Redirect(w, r, auth.AuthURL(state), 302)
	return
}

func loginCallbackHandler(w http.ResponseWriter, r *http.Request) {

	vars := structs.TemplateVars{}

	// Check if logged in
	loggedIn, err := session.IsLoggedIn(r)
	if err != nil {

		vars.ErrorMessage = "Spotify: " + err.Error()

		returnTemplate(w, r, "error", vars, err)
		return
	}
	if loggedIn {
		postlogin(w, r)
		return
	}

	// Check for OAuth errors
	queryErr := r.URL.Query().Get("error")
	if queryErr != "" {

		vars.ErrorMessage = "Spotify: " + queryErr

		returnTemplate(w, r, "error", vars, nil)
		return
	}

	auth := spotify.GetAuthenticator(r)
	state, err := session.Read(r, session.AuthState)
	if err != nil {
		rollbar.ErrorError(err)
	}

	// Check state and get token
	tok, err := auth.Token(state, r)
	if err != nil {

		vars.ErrorMessage = "Spotify: " + err.Error()

		returnTemplate(w, r, "error", vars, err)
		return
	}

	// Save user info to cookie
	client := auth.NewClient(tok)
	user, err := client.CurrentUser()
	if err != nil {

		vars.ErrorMessage = "Spotify: " + err.Error()

		returnTemplate(w, r, "error", vars, err)
		return
	}

	session.WriteMany(w, r, map[string]string{
		session.AuthState:    "",
		session.UserCountry:  user.Country,
		session.UserID:       user.ID,
		session.TokenToken:   tok.AccessToken,
		session.TokenType:    tok.TokenType,
		session.TokenRefresh: tok.RefreshToken,
		session.TokenExpiry:  strconv.Itoa(int(tok.Expiry.Unix())),
	})

	postlogin(w, r)
	return
}

func postlogin(w http.ResponseWriter, r *http.Request) {

	session.SetFlash(w, r, "Logged In :)")
	lastPage, err := session.Read(r, session.LastPage)
	if err != nil {
		rollbar.ErrorError(err)
	}

	http.Redirect(w, r, lastPage, 302)
	return
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {

	session.Clear(w, r)
	session.SetFlash(w, r, "Logged Out!")

	http.Redirect(w, r, "/", 302)
}
