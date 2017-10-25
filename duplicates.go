package main

import (
	"net/http"

	"github.com/Jleagle/spotifyhelper/session"
)

func duplicatesHandler(w http.ResponseWriter, r *http.Request) {

	session.Write(w, r, session.LastPage, "/duplicates")

	returnTemplate(w, r, "duplicates", templateVars{}, nil)
	return
}
