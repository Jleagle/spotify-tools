package main

import (
	"net/http"

	"github.com/Jleagle/spotify-tools/session"
)

func shuffleHandler(w http.ResponseWriter, r *http.Request) {

	//client := auth.NewClient(tok)

	session.Write(w, r, map[string]string{
		session.LastPage: "shuffle",
	})

	templateVars := templateVars{}
	templateVars.LoggedIn = session.IsLoggedIn(w, r)
	templateVars.Flashes = session.GetFlashes(w, r)

	returnTemplate(w, "shuffle", templateVars)
	return
}
