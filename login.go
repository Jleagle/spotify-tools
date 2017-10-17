package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Jleagle/spotify-tools/session"
	"github.com/zmb3/spotify"
)

var (
	auth  = spotify.NewAuthenticator("http://localhost:8084/login-callback", spotify.ScopeUserReadPrivate)
	state = "abc123"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {

	if session.IsLoggedIn(w, r) {
		postlogin(w, r)
	} else {
		auth.SetAuthInfo(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
		http.Redirect(w, r, auth.AuthURL(state), 302)
	}
}

func loginCallbackHandler(w http.ResponseWriter, r *http.Request) {

	if !session.IsLoggedIn(w, r) {

		auth.SetAuthInfo(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))

		tok, err := auth.Token(state, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusForbidden)
			log.Fatal(err)
		}
		if st := r.FormValue("state"); st != state {
			http.NotFound(w, r)
			log.Fatalf("State mismatch: %s != %s\n", st, state)
		}

		// Save to session
		session.Write(w, r, map[string]string{
			session.Token: tok.AccessToken,
		})
	}

	postlogin(w, r)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {

	session.Write(w, r, map[string]string{
		session.Token: "",
	})

	session.SetFlash(w, r, "Logged Out!")

	http.Redirect(w, r, "/", 302)
}

func postlogin(w http.ResponseWriter, r *http.Request) {

	session.SetFlash(w, r, "Logged In :)")
	lastPage := session.Read(w, r, session.LastPage)
	http.Redirect(w, r, "/"+lastPage, 302)
}
