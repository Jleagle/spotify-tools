package main

import (
	"net/http"

	"github.com/Jleagle/spotify-tools/session"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {

	session.Write(w, r, map[string]string{
		session.LastPage: "",
	})

	templateVars := templateVars{}
	templateVars.LoggedIn = session.IsLoggedIn(w, r)
	templateVars.Flashes = session.GetFlashes(w, r)

	returnTemplate(w, "home", templateVars)
	return
}
