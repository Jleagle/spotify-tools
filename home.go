package main

import (
	"net/http"

	"github.com/Jleagle/spotifyhelper/session"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {

	session.Write(w, r, session.LastPage, "/")

	returnTemplate(w, r, "home", templateVars{}, nil)
	return
}
