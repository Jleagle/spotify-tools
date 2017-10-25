package main

import (
	"net/http"

	"github.com/Jleagle/spotifyhelper/session"
)

func infoHandler(w http.ResponseWriter, r *http.Request) {

	session.Write(w, r, session.LastPage, "/info")

	returnTemplate(w, r, "info", templateVars{}, nil)
	return
}
