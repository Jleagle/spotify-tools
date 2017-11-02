package main

import (
	"fmt"
	"net/http"
	"strconv"

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

	state := helpers.RandomString(6, "abcdefghijklmnopqrstuvwxyz")
	session.Write(w, r, session.AuthState, state)

	auth := spot.GetAuthenticator(r)
	http.Redirect(w, r, auth.AuthURL(state), 302)
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

	auth := spot.GetAuthenticator(r)
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
	lastPage := session.Read(r, session.LastPage)
	http.Redirect(w, r, lastPage, 302)
	return
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {

	lastPage := session.Read(r, session.LastPage)

	session.Clear(w, r)
	session.SetFlash(w, r, "Logged Out!")

	http.Redirect(w, r, lastPage, 302)
}
